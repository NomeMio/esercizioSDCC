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
	server := rpc.NewServer() // create a server
	err := server.Register(handler)
	utilis.CheckError(err)
	lis, err := net.Listen("tcp", addr) // create a listener that handles RPCs
	defer lis.Close()
	utilis.CheckError(err)
	log.Printf("RPC server listens on port %s", *port)
	// Go specs: The caller typically invokes Accept in a go statement
	go func() {
		server.Accept(lis) // register the listener and accept inbound RPCs
	}()
	select {
	case <-handler.DoneWorking:
		return
	}
}
