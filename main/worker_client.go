package main

import (
	localrpc "esercizioSDCC/rpc"
	"flag"
	"net"
	"net/rpc"
	//"fmt"
	"esercizioSDCC/utilis"
	"log"
)

func main() {
	port := flag.String("port", "", "The port to connect to")
	flag.Parse()
	if *port == "" {
		log.Fatal("Must specify a port")
	}
	handler := new(localrpc.MapRequest)
	server := rpc.NewServer() // create a server
	err := server.Register(handler)
	utilis.CheckError(err)
	addr := "localhost:" + *port
	lis, err := net.Listen("tcp", addr) // create a listener that handles RPCs
	utilis.CheckError(err)
	log.Printf("RPC server listens on port %d", 8888)
	// Go specs: The caller typically invokes Accept in a go statement
	go func() {
		for {
			server.Accept(lis) // register the listener and accept inbound RPCs
		}
	}()
	select {}
}
