package main

import (
	"flag"
	"net"
	"net/rpc"
	localrpc "worker/rpc"
	"worker/utilis"

	"log"
)

func main() {
	port := flag.String("p", "", "The port to connect to")
	flag.Parse()
	if *port == "" {
		log.Fatal("Must specify a port")
	}
	addr := ":" + *port

	handler := localrpc.NewMapRequest(addr)
	server := rpc.NewServer()
	err := server.Register(handler)
	utilis.CheckError(err)
	lis, err := net.Listen("tcp", addr)
	defer lis.Close()
	utilis.CheckError(err)
	log.Printf("RPC server listens on port %s", *port)
	go func() {
		server.Accept(lis) // Attendi la chiamata del master
	}()
	select {
	case <-handler.DoneWorking:
		return
	}
}
