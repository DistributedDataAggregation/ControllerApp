package main

import (
	"encoding/binary"
	"log"
	"net"
	"strings"
	"sync"

	"controller/protomodels"

	"google.golang.org/protobuf/proto"
)

func sendToExecutors(files []string, queryReq HttpQueryRequest) error {

	filesPerExecutor, executors := distributeFiles(files, config.ExecutorAddresses)
	mainExecutor, otherExecutors := selectMainExecutor(executors)

	conn, err := sendTaskToExecutor(filesPerExecutor[mainExecutor], mainExecutor, mainExecutor, int32(len(executors)), queryReq, nil)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, executor := range otherExecutors {
		wg.Add(1)
		go sendTaskToExecutor(filesPerExecutor[executor], executor, mainExecutor, int32(len(executors)), queryReq, &wg)
	}
	wg.Wait()

	response, err := receiveResponseFromExecutor(conn)
	if err != nil {
		return err
	}

	log.Printf("Response from %s: %s", conn.RemoteAddr(), response)

	return nil
}

func selectMainExecutor(executors []string) (string, []string) {
	mainExecutor := executors[0]
	otherExecutors := executors[1:]
	return mainExecutor, otherExecutors
}

func distributeFiles(files []string, executors []string) (map[string][]string, []string) {
	filesPerExecutor := make(map[string][]string)
	usedExecutors := []string{}

	for i, file := range files {
		executor := executors[i%len(executors)]
		filesPerExecutor[executor] = append(filesPerExecutor[executor], file)

		if len(filesPerExecutor[executor]) == 1 {
			usedExecutors = append(usedExecutors, executor)
		}
	}
	return filesPerExecutor, usedExecutors
}

func createProtoRequest(files []string, queryReq HttpQueryRequest, mainExecutor string, isCurrentNodeMain bool, executorsCount int32) *protomodels.QueryRequest {
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
		Executor: &protomodels.MainExecutor{
			IsCurrentNodeMain: isCurrentNodeMain,
			IpAddress:         strings.Split(mainExecutor, ":")[0],
			Port:              8081,
			ExecutorsCount:    executorsCount,
		},
	}
}

func sendTaskToExecutor(files []string, executor string, mainExecutor string, executorsCount int32, queryReq HttpQueryRequest, wg *sync.WaitGroup) (net.Conn, error) {

	if wg != nil {
		defer wg.Done()
	}

	queryRequest := createProtoRequest(files, queryReq, mainExecutor, false, executorsCount)

	conn, err := sendRequest(queryRequest, executor, executor == mainExecutor)

	if err != nil {
		log.Printf("Error sending request to executor %s: %v", executor, err)
		return nil, err
	}

	return conn, nil
}

func sendRequest(queryRequest *protomodels.QueryRequest, executor string, isMain bool) (net.Conn, error) {
	conn, err := net.Dial("tcp", executor)
	if err != nil {
		return nil, err
	}

	size := proto.Size(queryRequest)
	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, uint32(size))

	data, err := proto.Marshal(queryRequest)
	if err != nil {
		conn.Close()
		log.Printf("Marshal error: %v", err)
		return nil, err
	}

	//message := append(sizeBytes, data...)
	_, err = conn.Write(sizeBytes)
	if err != nil {
		conn.Close()
		log.Printf("Error writing data to connection with %s: %v", executor, err)
		return nil, err
	}

	_, err = conn.Write(data)
	if err != nil {
		conn.Close()
		log.Printf("Error writing data to connection with %s: %v", executor, err)
		return nil, err
	}

	if isMain {
		return conn, nil
	}

	err = conn.Close()
	return nil, err
}

func receiveResponseFromExecutor(conn net.Conn) (string, error) {
	defer conn.Close()

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Printf("Error reading data from connection with main executor: %v", err)
		return "", err
	}

	return string(buffer[:n]), nil
}
