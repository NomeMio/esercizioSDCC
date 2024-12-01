package main

import (
	"esercizioSDCC/configuration"
	localrpc "esercizioSDCC/rpc"
	"esercizioSDCC/utilis"
	"fmt"
	"log"
	"net/rpc"
)

func makeRequestSync(client *rpc.Client, input string) ([2]int, error) {
	args := localrpc.MapArguemnt{InputString: input} // create the RPC arguments
	reply := localrpc.MapReply{}
	log.Printf("Synchronous call to RPC server")
	err := client.Call("MapRequest.MapGetResult", args, &reply)
	if err != nil {
		return [2]int{}, err
	}

	return [2]int{reply.MinValue, reply.MaxValue}, nil
}

func startWorker(addr string, argument string) error {
	client, err := rpc.Dial("tcp", addr)
	defer client.Close()
	utilis.CheckError(err)
	wc1, err1 := makeRequestSync(client, argument)
	utilis.CheckError(err1)
	println(addr, "worker started")
	fmt.Printf("Result: %v\n\n", wc1)
	return nil
}

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
	for key, token := range array {
		go func() {
			err := startWorker(configuration.HOSTS[key], token)
			utilis.CheckError(err)
		}()
	}
	select {}
}
