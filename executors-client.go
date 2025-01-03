package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"controller/protomodels"

	"google.golang.org/protobuf/proto"
)

const maxRetries = 10

type ExecutorsClient struct {
	MainIdx        int
	Sockets        []net.Conn
	Addresses      []string
	SocketStatuses []bool
	Mutex          *sync.Mutex
}

func NewExecutorsClient() *ExecutorsClient {
	sockets, addresses := OpenSockets(config.ExecutorAddresses, config.MainExecutorIdx)
	statuses := make([]bool, len(addresses))
	for i := range statuses {
		statuses[i] = sockets[i] != nil
	}

	return &ExecutorsClient{
		MainIdx:        0,
		Sockets:        sockets,
		Addresses:      addresses,
		SocketStatuses: statuses,
		Mutex:          &sync.Mutex{},
	}
}

func OpenSockets(executors []string, mainIdx int) ([]net.Conn, []string) {
	sockets := []net.Conn{}
	addresses := []string{}
	sockets = append(sockets, OpenSocket(executors[mainIdx]))
	addresses = append(addresses, executors[mainIdx])
	for i, executor := range executors {
		if i != mainIdx {
			sockets = append(sockets, OpenSocket(executor))
			addresses = append(addresses, executor)
		}
	}
	return sockets, addresses
}

func (ec *ExecutorsClient) ReconnectExecutor(executorIdx int) {
	ec.Mutex.Lock()
	defer ec.Mutex.Unlock()

	if !ec.SocketStatuses[executorIdx] {
		conn := OpenSocket(ec.Addresses[executorIdx])
		ec.Sockets[executorIdx] = conn
		ec.SocketStatuses[executorIdx] = true
		log.Printf("Reconnected to executor %s", ec.Addresses[executorIdx])
	}
}

func OpenSocket(executor string) net.Conn {

	var conn net.Conn
	var err error

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
	return conn
}

func (ec *ExecutorsClient) allExecutorsConnected() bool {
	for _, status := range ec.SocketStatuses {
		if !status {
			return false
		}
	}
	return true
}

func (ec *ExecutorsClient) createProtoRequest(guid string, files []string, queryReq HttpQueryRequest, mainExecutor string, mainExecutorPort int32,
	isCurrentNodeMain bool, executorsCount int32) (*protomodels.QueryRequest, error) {
	selects := make([]*protomodels.Select, len(queryReq.SelectColumns))
	for i, sel := range queryReq.SelectColumns {

		if !sel.Function.IsValid() {
			log.Printf("Invalid aggreage function %s. Supported aggregate functions: Minimum, Maximum, Average", string(sel.Function))
			return nil, fmt.Errorf("invalid aggreage function %s, supported aggregate functions: Minimum, Maximum, Average", string(sel.Function))
		}

		selects[i] = &protomodels.Select{
			Column:   sel.Column,
			Function: protomodels.Aggregate(protomodels.Aggregate_value[string(sel.Function)]),
		}
	}

	return &protomodels.QueryRequest{
		Guid:         guid,
		FilesNames:   files,
		GroupColumns: queryReq.GroupColumns,
		Select:       selects,
		Executor: &protomodels.ExecutorInformation{
			IsCurrentNodeMain: isCurrentNodeMain,
			MainIpAddress:     strings.Split(mainExecutor, ":")[0],
			MainPort:          mainExecutorPort,
			ExecutorsCount:    executorsCount,
		},
	}, nil
}

func (ec *ExecutorsClient) sendTaskToExecutor(guid string, files []string, executorIdx int, executorsCount int32, queryReq HttpQueryRequest) error {

	queryRequest, err := ec.createProtoRequest(guid, files, queryReq, ec.Addresses[ec.MainIdx], config.ExecutorsPort, executorIdx == ec.MainIdx, executorsCount)
	if err != nil {
		return err
	}

	err = ec.sendRequest(queryRequest, executorIdx)
	if err != nil {
		log.Printf("Error sending request to executor %s: %v", ec.Sockets[executorIdx], err)
		return err
	}

	return nil
}

func (ec *ExecutorsClient) printProtoRequest(queryReq *protomodels.QueryRequest, adress net.Addr) {
	log.Printf("Sent request to %s\n", adress)
	log.Printf("Files:\n")
	for _, file := range queryReq.FilesNames {
		log.Printf("	%s, \n", file)
	}

	log.Printf("Grouping columns:\n")
	for _, col := range queryReq.Select {
		log.Printf("	%s, \n", col)
	}

	log.Printf("Select:\n")
	for _, sel := range queryReq.Select {
		log.Printf("	Column %s, Function %s,\n", sel.Column, sel.Function)
	}

	log.Printf("Main executor:\n")
	log.Printf("	Is main: %t", queryReq.Executor.IsCurrentNodeMain)
	log.Printf("	Main ip address: %s", queryReq.Executor.MainIpAddress)
	log.Printf("	Main port: %d", queryReq.Executor.MainPort)
	log.Printf("	Executors count: %d", queryReq.Executor.ExecutorsCount)

}

func (ec *ExecutorsClient) sendRequest(queryRequest *protomodels.QueryRequest, executorIdx int) error {

	size := proto.Size(queryRequest)
	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, uint32(size))

	data, err := proto.Marshal(queryRequest)
	if err != nil {
		log.Printf("Marshal error: %v", err)
		return err
	}

	err = ec.connWrite(sizeBytes, executorIdx)
	if err != nil {
		return err
	}

	err = ec.connWrite(data, executorIdx)

	return err
}

func (ec *ExecutorsClient) connWrite(data []byte, executorIdx int) error {
	conn := ec.Sockets[executorIdx]
	_, err := conn.Write(data)
	if err != nil {
		conn.Close()
		log.Printf("Error writing data to connection with %s: %v", ec.Addresses[executorIdx], err)
		ec.Mutex.Lock()
		ec.SocketStatuses[executorIdx] = false
		ec.Mutex.Unlock()
		go ec.ReconnectExecutor(executorIdx)
		return err
	}
	return nil
}

func (ec *ExecutorsClient) receiveResponseFromMainExecutor(guid string) (HttpResult, error) {

	for i := 0; i < maxRetries; i++ {

		sizeBytes, err := ec.connRead(4, ec.MainIdx)
		if err != nil {
			return HttpResult{}, err
		}
		messageSize := binary.BigEndian.Uint32(sizeBytes)

		data, err := ec.connRead(int(messageSize), ec.MainIdx)
		if err != nil {
			return HttpResult{}, err
		}

		response, receivedGuid, err := ec.readResponseFromMainExecutor(data)

		if err != nil {
			return HttpResult{}, err
		}

		if guid == receivedGuid {
			return response, nil
		}

	}

	return HttpResult{}, fmt.Errorf("failed to receive response from %s after %d retries", ec.Addresses[ec.MainIdx], maxRetries)
}

func (ec *ExecutorsClient) connRead(size int, executorIdx int) ([]byte, error) {
	conn := ec.Sockets[executorIdx]
	data := make([]byte, size)
	_, err := io.ReadFull(conn, data)
	if err != nil {
		conn.Close()
		log.Printf("Error reading data from connection with %s: %v", ec.Addresses[executorIdx], err)
		ec.Mutex.Lock()
		ec.SocketStatuses[executorIdx] = false
		ec.Mutex.Unlock()
		go ec.ReconnectExecutor(executorIdx)
		return []byte{}, err
	}
	return data, nil
}

func (ec *ExecutorsClient) readResponseFromMainExecutor(data []byte) (HttpResult, string, error) {

	if len(data) == 0 {
		return HttpResult{}, "", fmt.Errorf("empty input data")
	}

	var queryResponse protomodels.QueryResponse
	err := proto.Unmarshal(data, &queryResponse)
	if err != nil {
		log.Printf("Error unmarshalling QueryResponse: %v", err)
		return HttpResult{}, "", err
	}

	httpResult := HttpResult{
		Response: ec.mapQueryResponse(&queryResponse),
	}

	return httpResult, queryResponse.Guid, nil
}

func (ec *ExecutorsClient) mapQueryResponse(src *protomodels.QueryResponse) HttpQueryResponse {
	if src == nil {
		return HttpQueryResponse{}
	}

	var httpError *HttpError
	if src.Error != nil {
		httpError = &HttpError{
			Message:      src.Error.Message,
			InnerMessage: src.Error.InnerMessage,
		}
	}

	httpValues := make([]*HttpValue, len(src.Values))
	for i, value := range src.Values {
		if value != nil {
			httpValues[i] = &HttpValue{
				GroupingValue: value.GroupingValue,
				Results:       ec.mapPartialResults(value.Results),
			}
		}
	}

	return HttpQueryResponse{
		Error:  httpError,
		Values: httpValues,
	}
}

func (ec *ExecutorsClient) mapPartialResults(results []*protomodels.PartialResult) []HttpPartialResult {
	httpResults := make([]HttpPartialResult, len(results))
	for i, result := range results {
		if result != nil {
			httpResults[i] = HttpPartialResult{
				Value: result.Value,
				Count: result.Count,
			}
		}
	}
	return httpResults
}
