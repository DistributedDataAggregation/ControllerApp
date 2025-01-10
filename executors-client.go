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
	MainIdx        *int
	MainIdxMutex   *sync.Mutex
	Sockets        []net.Conn
	Addresses      []string
	SocketStatuses []bool
	Mutexes        []*sync.Mutex
}

func NewExecutorsClient() *ExecutorsClient {

	mainIdx := config.MainExecutorIdx

	addresses := []string{}
	addresses = append(addresses, config.ExecutorAddresses[mainIdx])
	for i, executor := range config.ExecutorAddresses {
		if i != mainIdx {
			addresses = append(addresses, executor)
		}
	}

	mainIdx = -1

	statuses := make([]bool, len(addresses))
	for i := range statuses {
		statuses[i] = false
	}

	socketMutexes := make([]*sync.Mutex, len(addresses))
	for i := range addresses {
		socketMutexes[i] = &sync.Mutex{}
	}

	return &ExecutorsClient{
		MainIdx:        &mainIdx,
		MainIdxMutex:   &sync.Mutex{},
		Sockets:        make([]net.Conn, len(addresses)),
		Addresses:      addresses,
		SocketStatuses: statuses,
		Mutexes:        socketMutexes,
	}
}

func (ec *ExecutorsClient) OpenSockets() {
	for i := range ec.Addresses {
		go ec.connectToExecutor(i)
	}
}

func (ec *ExecutorsClient) connectToExecutor(executorIdx int) {

	ec.Mutexes[executorIdx].Lock()
	defer ec.Mutexes[executorIdx].Unlock()

	log.Printf("Connecting to executor %s", ec.Addresses[executorIdx])

	for {
		conn, err := ec.openSocket(ec.Addresses[executorIdx])
		if err == nil {
			ec.Sockets[executorIdx] = conn
			ec.SocketStatuses[executorIdx] = true
			log.Printf("Connected to executor %s", ec.Addresses[executorIdx])

			if executorIdx == 0 || *ec.MainIdx == -1 {
				ec.MainIdxMutex.Lock()
				*ec.MainIdx = executorIdx
				ec.MainIdxMutex.Unlock()
			}
			return
		}
		log.Printf("Retrying connection to executor %s", ec.Addresses[executorIdx])
		time.Sleep(2 * time.Second)
	}
}

func (ec *ExecutorsClient) openSocket(executor string) (net.Conn, error) {
	var conn net.Conn
	var err error

	baseDelay := time.Second

	for retries := 0; retries < maxRetries; retries++ {
		conn, err = net.Dial("tcp", executor)
		if err == nil {
			return conn, nil
		}

		backoff := baseDelay * (1 << retries)
		log.Printf("Retrying to connect to %v after %v... (%d/%d)", executor, backoff, retries+1, maxRetries)
		time.Sleep(backoff)
	}

	if err != nil {
		log.Printf("Failed to connect to %v after %d retries", executor, maxRetries)
	}

	return nil, fmt.Errorf("failed to connect to %v after %d retries", executor, maxRetries)
}

func (ec *ExecutorsClient) GetAvailableExecutorIdxs() ([]int, error) {

	if *ec.MainIdx == -1 || !ec.SocketStatuses[*ec.MainIdx] {
		return nil, fmt.Errorf("main executor is unavailable")
	}

	available := []int{*ec.MainIdx}
	for i, status := range ec.SocketStatuses {
		if i != *ec.MainIdx && status {
			available = append(available, i)
		}
	}
	return available, nil
}

func (ec *ExecutorsClient) getFirstAvailableExecutor() int {
	for i, status := range ec.SocketStatuses {
		if status {
			return i
		}
	}
	return -1
}

func (ec *ExecutorsClient) SendTaskToExecutor(guid string, files []string, executorIdx int, executorsCount int32, queryReq HttpQueryRequest) error {

	//ports := []int32{8081, 8083, 8085}
	//queryRequest, err := CreateProtoRequest(guid, files, queryReq, ec.Addresses[*ec.MainIdx], ports[*ec.MainIdx], executorIdx == *ec.MainIdx, executorsCount)

	queryRequest, err := CreateProtoRequest(guid, files, queryReq, ec.Addresses[*ec.MainIdx], config.ExecutorsPort, executorIdx == *ec.MainIdx, executorsCount)
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
		ec.deactivateConn(executorIdx)
		return err
	}
	return nil
}

func (ec *ExecutorsClient) ReceiveResponseFromMainExecutor(guid string) (HttpResult, error) {

	for i := 0; i < maxRetries; i++ {

		sizeBytes, err := ec.connRead(4, *ec.MainIdx)
		if err != nil {
			return HttpResult{}, err
		}
		messageSize := binary.BigEndian.Uint32(sizeBytes)

		data, err := ec.connRead(int(messageSize), *ec.MainIdx)
		if err != nil {
			return HttpResult{}, err
		}

		response, receivedGuid, err := ReadQueryResultProto(data)

		if err != nil {
			return HttpResult{}, err
		}

		if guid == receivedGuid {
			return response, nil
		}

	}

	return HttpResult{}, fmt.Errorf("failed to receive response from %s after %d retries", ec.Addresses[*ec.MainIdx], maxRetries)
}

func (ec *ExecutorsClient) connRead(size int, executorIdx int) ([]byte, error) {
	conn := ec.Sockets[executorIdx]
	data := make([]byte, size)
	_, err := io.ReadFull(conn, data)
	if err != nil {
		conn.Close()
		log.Printf("Error reading data from connection with %s: %v", ec.Addresses[executorIdx], err)
		ec.deactivateConn(executorIdx)
		return []byte{}, err
	}
	return data, nil
}

func (ec *ExecutorsClient) deactivateConn(executorIdx int) {

	ec.Mutexes[executorIdx].Lock()
	ec.SocketStatuses[executorIdx] = false
	ec.Mutexes[executorIdx].Unlock()

	if executorIdx == *ec.MainIdx {

		ec.MainIdxMutex.Lock()
		newMainIdx := ec.getFirstAvailableExecutor()
		*ec.MainIdx = newMainIdx
		ec.MainIdxMutex.Unlock()

		go ec.connectToExecutor(executorIdx)
	}
}
