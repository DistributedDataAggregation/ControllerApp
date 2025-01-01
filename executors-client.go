package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"controller/protomodels"

	"google.golang.org/protobuf/proto"
)

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
	return conn
}

func (ec *ExecutorsClient) AllExecutorsConnected() bool {
	ec.Mutex.Lock()
	defer ec.Mutex.Unlock()

	for _, status := range ec.SocketStatuses {
		if !status {
			return false
		}
	}
	return true
}

func (ec *ExecutorsClient) createProtoRequest(files []string, queryReq HttpQueryRequest, mainExecutor string, mainExecutorPort int32, isCurrentNodeMain bool, executorsCount int32) *protomodels.QueryRequest {
	selects := make([]*protomodels.Select, len(queryReq.SelectColumns))
	for i, sel := range queryReq.SelectColumns {
		selects[i] = &protomodels.Select{
			Column:   sel.Column,
			Function: protomodels.Aggregate(protomodels.Aggregate_value[sel.Function]), // TODO handle unsupported aggregate (now it returns default Minimum)
		}
	}

	return &protomodels.QueryRequest{
		FilesNames:   files,
		GroupColumns: queryReq.GroupColumns,
		Select:       selects,
		Executor: &protomodels.ExecutorInformation{
			IsCurrentNodeMain: isCurrentNodeMain,
			MainIpAddress:     strings.Split(mainExecutor, ":")[0],
			MainPort:          mainExecutorPort,
			ExecutorsCount:    executorsCount,
		},
	}
}

func (ec *ExecutorsClient) sendTaskToExecutor(files []string, executorIdx int, executorsCount int32, queryReq HttpQueryRequest, wg *sync.WaitGroup) error {

	if wg != nil {
		defer wg.Done()
	}

	queryRequest := ec.createProtoRequest(files, queryReq, ec.Addresses[ec.MainIdx], config.ExecutorsPort, executorIdx == ec.MainIdx, executorsCount)

	err := ec.sendRequest(queryRequest, executorIdx)

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

func (ec *ExecutorsClient) receiveResponseFromMainExecutor() (HttpResult, error) {

	sizeBytes := make([]byte, 4)
	err := ec.connRead(&sizeBytes, ec.MainIdx)
	if err != nil {

		return HttpResult{}, err
	}
	messageSize := binary.BigEndian.Uint32(sizeBytes)

	data := make([]byte, messageSize)
	err = ec.connRead(&data, ec.MainIdx)
	if err != nil {

		return HttpResult{}, err
	}

	return ec.readResponseFromMainExecutor(data)
}

func (ec *ExecutorsClient) connRead(data *[]byte, executorIdx int) error {
	conn := ec.Sockets[executorIdx]
	_, err := conn.Read(*data)
	if err != nil {
		conn.Close()
		log.Printf("Error reading data from connection with %s: %v", ec.Addresses[executorIdx], err)
		ec.Mutex.Lock()
		ec.SocketStatuses[executorIdx] = false
		ec.Mutex.Unlock()
		go ec.ReconnectExecutor(executorIdx)
		return err
	}
	return nil
}

func (ec *ExecutorsClient) readResponseFromMainExecutor(data []byte) (HttpResult, error) {

	if len(data) == 0 {
		return HttpResult{}, fmt.Errorf("empty input data")
	}

	var queryResponse protomodels.QueryResponse
	err := proto.Unmarshal(data, &queryResponse)
	if err != nil {
		log.Printf("Error unmarshalling QueryResult: %v", err)
		return HttpResult{}, err
	}

	httpResult := HttpResult{
		Response: ec.mapQueryResponse(&queryResponse),
	}

	return httpResult, nil
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
