package main

import (
	"encoding/binary"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"controller/protomodels"

	"google.golang.org/protobuf/proto"
)

type ExecutorsClient struct {
	MainIdx int
	Sockets []net.Conn
}

func NewExecutorsClient() *ExecutorsClient {
	return &ExecutorsClient{MainIdx: 0, Sockets: OpenSockets(config.ExecutorAddresses)}
}

func OpenSockets(executors []string) []net.Conn {
	sockets := []net.Conn{}
	for _, executor := range executors {
		var conn net.Conn
		var err error
		const maxRetries = 10
		baseDelay := time.Second // Base delay for retries

		for retries := 0; retries < maxRetries; retries++ {
			conn, err = net.Dial("tcp", executor)
			if err == nil {
				break
			}

			// Calculate exponential backoff
			backoff := baseDelay * (1 << retries) // 1, 2, 4, 8, 16 seconds
			log.Printf("Retrying to connect to %v after %v... (%d/%d)", executor, backoff, retries+1, maxRetries)
			time.Sleep(backoff)
		}

		if err != nil {
			log.Panicf("Failed to dial connect to %v after %d retries: %v", executor, maxRetries, err)
		}
		sockets = append(sockets, conn)
	}
	return sockets
}

func (ec *ExecutorsClient) createProtoRequest(files []string, queryReq HttpQueryRequest, mainExecutor string, isCurrentNodeMain bool, executorsCount int32) *protomodels.QueryRequest {
	selects := make([]*protomodels.Select, len(queryReq.SelectColumns))
	for i, sel := range queryReq.SelectColumns {
		selects[i] = &protomodels.Select{
			Column:   sel.Column,
			Function: protomodels.Aggregate(protomodels.Aggregate_value[strings.ToUpper(sel.Function)]),
		}
	}

	return &protomodels.QueryRequest{
		FilesNames:   files,
		GroupColumns: queryReq.GroupColumns,
		Select:       selects,
		Executor: &protomodels.ExecutorInformation{
			IsCurrentNodeMain: isCurrentNodeMain,
			MainIpAddress:     "172.20.0.2",
			MainPort:          8081, // TODO int parse strings.Split(mainExecutor, ":")[1],
			ExecutorsCount:    executorsCount,
		},
	}
}

func (ec *ExecutorsClient) sendTaskToExecutor(files []string, executorIdx int, executorsCount int32, queryReq HttpQueryRequest, wg *sync.WaitGroup) error {

	if wg != nil {
		defer wg.Done()
	}

	queryRequest := ec.createProtoRequest(files, queryReq, config.ExecutorAddresses[ec.MainIdx], executorIdx == ec.MainIdx, executorsCount)

	err := ec.sendRequest(queryRequest, ec.Sockets[executorIdx])

	if err != nil {
		log.Printf("Error sending request to executor %s: %v", ec.Sockets[executorIdx], err)
		return err
	}

	return nil
}

func (ec *ExecutorsClient) sendRequest(queryRequest *protomodels.QueryRequest, conn net.Conn) error {

	size := proto.Size(queryRequest)
	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, uint32(size))

	data, err := proto.Marshal(queryRequest)
	if err != nil {
		conn.Close()
		log.Printf("Marshal error: %v", err)
		return err
	}

	_, err = conn.Write(sizeBytes)
	if err != nil {
		conn.Close()
		log.Printf("Error writing data to connection with %s: %v", conn.RemoteAddr(), err)
		return err
	}

	_, err = conn.Write(data)
	if err != nil {
		conn.Close()
		log.Printf("Error writing data to connection with %s: %v", conn.RemoteAddr(), err)
		return err
	}

	return nil
}

func (ec *ExecutorsClient) receiveResponseFromMainExecutor() (HttpResult, error) {
	conn := ec.Sockets[ec.MainIdx]

	sizeBytes := make([]byte, 4)
	_, err := conn.Read(sizeBytes)
	if err != nil {
		log.Printf("Error reading size from connection with main executor: %v", err)
		return HttpResult{}, err
	}
	messageSize := binary.BigEndian.Uint32(sizeBytes)

	data := make([]byte, messageSize)
	_, err = conn.Read(data)
	if err != nil {
		log.Printf("Error reading data from connection with main executor: %v", err)
		return HttpResult{}, err
	}

	return ec.readResponseFromMainExecutor(data)
}

func (ec *ExecutorsClient) readResponseFromMainExecutor(data []byte) (HttpResult, error) {

	var queryResponse protomodels.QueryResponse
	err := proto.Unmarshal(data, &queryResponse)
	if err != nil {
		log.Printf("Error unmarshalling QueryResult: %v", err)
		return HttpResult{}, err
	}

	httpResult := HttpResult{
		Response: queryResponse,
	}

	return httpResult, nil
}
