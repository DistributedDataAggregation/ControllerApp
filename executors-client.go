package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
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

func (ec *ExecutorsClient) sendTaskToExecutor(guid string, files []string, executorIdx int, executorsCount int32, queryReq HttpQueryRequest) error {

	queryRequest, err := createProtoRequest(guid, files, queryReq, ec.Addresses[ec.MainIdx], config.ExecutorsPort, executorIdx == ec.MainIdx, executorsCount)
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

		response, receivedGuid, err := readProtoResponse(data)

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
