package main

import (
	"esercizioSDCC/configuration"
	localrpc "esercizioSDCC/rpc"
	"esercizioSDCC/utilis"
	"flag"
	"fmt"
	"log"
	"math"
	"net/rpc"
	"regexp"
	"strings"
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
var hosts []string     //indirizzi degli host
var numbersOfHosts int //numero di host che ho = len(hosts)
var maxNUmber int
var numberOfInts int

func main() {
	// Inizializzo i parametri acquisiti come argomenti del programma
	hostsFlag := flag.String("a", "", "The workesrs addreses in format host:port,host2:port2,....") //Rappresenta gli indirizzi dei worker che verranno chiamati
	numberOfIntsGenerated := flag.Int("n", 1000, "Number of workers to use")                        //Rappresenta il numero di interi generati
	maxIntNumber := flag.Int("m", 1000, "Maximum number of workers to use")                         //Rappresenta il numero massimo generabili durante la creazione del file

	flag.Parse()

	if *hostsFlag == "" {
		flag.Usage()
		log.Fatal("Number of hosts must be greater than zero")
	}

	//Controlla che gli indirizzi siano legali attravesro una regex
	maxNUmber = *maxIntNumber
	numbersOfHosts = *numberOfIntsGenerated
	hostsSplitted := strings.Split(*hostsFlag, ",")
	regex := regexp.MustCompile(configuration.ADDREDSSPATTERN)
	for _, temp := range hostsSplitted {
		if !regex.MatchString(temp) {
			print("Invalid host: " + temp)
			return
		}
	}
	hosts = make([]string, len(hostsSplitted))
	copy(hosts, hostsSplitted)
	numbersOfHosts := len(hosts)

	fmt.Println("Server started")
	nome, err := utilis.GenerateRandomIntsFIle(numberOfInts, maxNUmber) //genero il file
	array := utilis.GetStrings(nome, numbersOfHosts, numberOfInts)      //splitto il file in shard per i client
	utilis.CheckError(err)
	fmt.Println("Generated ", numberOfIntsGenerated, "integers values at ", nome)
	starWorkers(array) //inizio a far lavorare gli worker

}

func starWorkers(array []string) {
	results := make([][2]int, numbersOfHosts)
	var wg sync.WaitGroup
	connections = make(map[int]connectionStruct)
	for key, token := range array {
		wg.Add(1)
		go func() {
			defer wg.Done()
			connections[key] = connectionStruct{startWorker(hosts[key]), hosts[key]}
			result, err := makeMapRequest(connections[key], token)
			utilis.CheckError(err)
			results[key] = result
			fmt.Printf("Resulting min and max from  mapping of hosy %s: %v\n", hosts[key], result)
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
	floatedNumber := float64(max-min) / float64(numbersOfHosts)
	numberOfKeys := int(math.Ceil(floatedNumber))
	shuffleChiavi := make([]localrpc.ReduceMap, numbersOfHosts)
	minIterator := 0
	maxIterator := min + numberOfKeys
	for number, host := range hosts {
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
