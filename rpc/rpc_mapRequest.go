package rpc

import (
	"encoding/json"
	"esercizioSDCC/utilis"
	"log"
	"net/rpc"
	"strconv"
	"strings"
	"sync"
)

// struttura su cui si registra del rpc
type MapRequest struct {
	frequency map[int][]int //valori dell'host
	lock      sync.Locker
	//aggiungi contatore host per sincronia
	address string
	done    chan bool
}

// argomento della chiamata MapGetResult che rappresenta lo shard del file
type MapArgument struct {
	InputString string
}

type ReduceArgument struct {
	InputMap map[int][]int
}

// ritorno della chiamata MapGetResult che ci indica il valore massimo e minimo tra le chiavi ricavate
type MapReply struct {
	MaxValue int
	MinValue int
}

// inizializza la struttura MapRequest da scrivere nel registro rpc
func NewMapRequest(host string) MapRequest {
	newObj := MapRequest{frequency: make(map[int][]int), lock: new(sync.Mutex)}
	newObj.address = host
	newObj.done = make(chan bool)
	return newObj
}

func (r MapRequest) MapAppending(key int, value int) {
	r.lock.Lock()
	r.frequency[key] = append(r.frequency[key], value)
	r.lock.Unlock()
}

// effettua la mappatura dello shard nella mappa frequency e ritorna il valore massim e minimo fra le chiavi della mappa
func (r MapRequest) MapGetResult(argument MapArgument, reply *MapReply) error {
	//TODO vedere cosa fare in caso di chiamate duplicate
	if len(r.frequency) != 0 {
		r.frequency = make(map[int][]int)
	}
	tokens := strings.Fields(argument.InputString)
	var minMax [2]int
	for value, token := range tokens {
		atoi, err := strconv.Atoi(token)
		utilis.CheckError(err)
		if value == 0 {
			minMax = [2]int{atoi, atoi}
		}
		if atoi < minMax[0] {
			minMax[0] = atoi
		} else if atoi > minMax[1] {
			minMax[1] = atoi
		}

		r.MapAppending(atoi, 1)
	}
	jsonString, err := json.Marshal(r.frequency)
	utilis.CheckError(err)
	print(string(jsonString))
	reply.MinValue = minMax[0]
	reply.MaxValue = minMax[1]
	return nil
}

func (r MapRequest) SendMapValues(argument ReduceArgument, reply *MapReply) error {
	// Ensure r.frequency is initialized
	if r.frequency == nil {
		r.frequency = make(map[int][]int) // Replace KeyType and ValueType with the actual types
	}

	// Add the values from 'argument' to 'r.frequency'
	for key, value := range argument.InputMap {
		for _, x := range value {
			r.MapAppending(key, x)
		}
	}
	r.done <- true
	return nil
}

func (r MapRequest) StartValuesExchange(arguments []ReduceMap, reply *ReduceReply) error {
	toWait := len(arguments) - 1
	for _, item := range arguments {
		if item.Host != r.address {
			//fmt.Printf("Host %s processing range: %v\n", r.address, item.KeyRange)

			// Filter data by key range
			filteredResults := make(map[int][]int)
			for key, value := range r.frequency { // Ciclo sui valori che il singolo host ha
				if key >= item.KeyRange[0] && key <= item.KeyRange[1] {
					filteredResults[key] = value // Mappato valore su nuovo range (corretto)
					delete(r.frequency, key)     // Eliminazione del valore mappato dal vecchio range
				}
			}
			// Print filtered results for debugging
			//fmt.Printf("Host %s filtered results: %v\n", r.address, filteredResults)
			argumentReduce := ReduceArgument{filteredResults}
			//Send Map Values
			go func() {
				client, err1 := rpc.Dial("tcp", item.Host) //inizializzo connesione
				utilis.CheckError(err1)
				log.Printf("Call to RPC server %s\n", item.Host)
				err2 := client.Call("MapRequest.SendMapValues", argumentReduce, nil)
				utilis.CheckError(err2)
			}()
		}
	}
	//todo aspetta contatore di tutti gli host
	return nil
}

type ReduceReply struct {
	Reply string
}

type ReduceMap struct {
	Host     string
	KeyRange [2]int
}
