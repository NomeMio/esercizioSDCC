package main

import (
	"esercizioSDCC/configuration"
	localrpc "esercizioSDCC/rpc"
	"esercizioSDCC/utilis"
	"fmt"
	"log"
	"math"
	"net/rpc"
	"sync"
)

type connectionStruct struct {
	conn *rpc.Client
	addr string
}

func makeMapRequest(client connectionStruct, input string) ([2]int, error) {
	args := localrpc.MapArguemnt{InputString: input} // create the RPC arguments
	reply := localrpc.MapReply{}
	log.Printf("Call to RPC server %s\n", client.addr)
	err := client.conn.Call("MapRequest.MapGetResult", args, &reply)
	if err != nil {
		return [2]int{}, err
	}

	return [2]int{reply.MinValue, reply.MaxValue}, nil
}

func makeReduceRequest(client connectionStruct, shuffledKeys []localrpc.ReduceMap) error {
	log.Printf("Call to RPC server %s\n", client.addr)
	reply := localrpc.ReduceReply{}
	err := client.conn.Call("MapRequest.ShowMyRange", shuffledKeys, &reply)
	utilis.CheckError(err)
	return nil
}

func startWorker(addr string) *rpc.Client {
	client, err := rpc.Dial("tcp", addr)
	utilis.CheckError(err)
	println(addr, "worker started")
	return client
}

var connections map[int]connectionStruct

func main() {

	fmt.Println("Server started")
	nome, err := utilis.GenerateRandomIntsFIle()
	array := utilis.GetStrings(nome)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Generated ", configuration.FILE_SIZE, "integers values at ", nome)
	starWorkers(array)

}

func starWorkers(array []string) {
	results := [configuration.HOSTS_NUMBER][2]int{}
	var wg sync.WaitGroup
	connections = make(map[int]connectionStruct)
	for key, token := range array {
		wg.Add(1)
		go func() {
			defer wg.Done()
			connections[key] = connectionStruct{startWorker(configuration.HOSTS[key]), configuration.HOSTS[key]}
			result, err := makeMapRequest(connections[key], token)
			utilis.CheckError(err)
			results[key] = result
			fmt.Printf("Resulting min and max from  mapping of hosy %s: %v\n", configuration.HOSTS[key], result)
		}()
	}
	wg.Wait()
	println("done mapping")
	var max, min int
	for i, result := range results {
		if i == 0 {
			min = result[0]
			max = result[1]
			continue
		}
		if max < result[1] {
			max = result[1]
		} else if min > result[0] {
			min = result[0]
		}
	}
	println("min:", min, "max:", max)
	floatedNumber := float64(max-min) / float64(configuration.HOSTS_NUMBER)
	numberOfKeys := int(math.Ceil(floatedNumber))
	shuffleChiavi := [configuration.HOSTS_NUMBER]localrpc.ReduceMap{}
	minIterator := 0
	maxIterator := min + numberOfKeys
	for number, host := range configuration.HOSTS {
		shuffleChiavi[number] = localrpc.ReduceMap{Host: host, KeyRange: [2]int{minIterator, maxIterator}}
		minIterator = maxIterator + 1
		maxIterator += numberOfKeys
	}
	fmt.Printf("map of keys:%v\n", shuffleChiavi)
	wg = sync.WaitGroup{}
	for _, token := range connections {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := makeReduceRequest(token, shuffleChiavi[:])
			utilis.CheckError(err)
			fmt.Printf("starting reduce in host %s\n", token.addr)
		}()
	}
	wg.Wait()

}
