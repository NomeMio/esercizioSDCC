package rpc

import (
	"esercizioSDCC/utilis"
	"net/rpc"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// struttura su cui si registra del rpc
type RpcMapReduce struct {
	frequency   map[int][]int //valori dell'host
	lock        sync.Locker   //aggiungi contatore host per sincronia
	address     string
	donePeers   chan bool
	DoneWorking chan bool
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

// inizializza la struttura RpcMapReduce da scrivere nel registro rpc
func NewMapRequest(host string) RpcMapReduce {
	newObj := RpcMapReduce{frequency: make(map[int][]int), lock: new(sync.Mutex)}
	newObj.address = host
	newObj.donePeers = make(chan bool)
	newObj.DoneWorking = make(chan bool)
	return newObj
}

func (r RpcMapReduce) MapAppending(key int, value int) {
	r.lock.Lock()
	r.frequency[key] = append(r.frequency[key], value)
	r.lock.Unlock()
}

func (r RpcMapReduce) MapRemove(key int) {
	r.lock.Lock()
	delete(r.frequency, key)
	r.lock.Unlock()
}

// effettua la mappatura dello shard nella mappa frequency e ritorna il valore massim e minimo fra le chiavi della mappa
func (r RpcMapReduce) MapGetResult(argument MapArgument, reply *MapReply) error {
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
	//jsonString, err := json.Marshal(r.frequency)
	//utilis.CheckError(err)
	//println("Ricevuto :", string(jsonString))
	reply.MinValue = minMax[0]
	reply.MaxValue = minMax[1]
	return nil
}

func (r RpcMapReduce) SendMapValues(argument ReduceArgument, reply *MapReply) error {
	// Add the values from 'argument' to 'r.frequency'
	for key, value := range argument.InputMap {
		for _, x := range value {
			r.MapAppending(key, x)
		}
	}
	r.donePeers <- true
	return nil
}

func (r RpcMapReduce) StartValuesExchange(arguments []ReduceMap, reply *ReduceReply) error {
	toWait := len(arguments) - 1
	for _, item := range arguments {
		if item.Host != r.address {
			// Filter data by key range
			filteredResults := make(map[int][]int)
			for key, value := range r.frequency { // Ciclo sui valori che il singolo host ha
				if key >= item.KeyRange[0] && key <= item.KeyRange[1] {
					filteredResults[key] = value // Mappato valore su nuovo range (corretto)
					r.MapRemove(key)             // Eliminazione del valore mappato dal vecchio range
				}
			}
			//log.Printf("Sending to %s values  %v\n", item.Host, filteredResults)
			argumentReduce := ReduceArgument{filteredResults}
			go func() {
				client, err1 := rpc.Dial("tcp", item.Host) //inizializzo connesione
				utilis.CheckError(err1)
				err2 := client.Call("RpcMapReduce.SendMapValues", argumentReduce, nil)
				utilis.CheckError(err2)
			}()
		}
	}
	for i := 1; i <= toWait; i++ {
		<-r.donePeers
	}
	//log.Printf("I'm the host %s con valori :\n%v\n", r.address, r.frequency)
	//log.Printf("Reducing\n")
	for key, value := range r.frequency {
		lenght := len(value)
		split := []int{lenght}
		r.frequency[key] = split
	}
	//log.Printf("%v\n", r.frequency)
	//log.Printf("%v\n", r.frequency)
	stringa := ""
	var keys []int
	for key := range r.frequency {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, key := range keys {
		value := r.frequency[key]
		for i := 0; i < value[0]; i++ {
			stringa += strconv.Itoa(key) + " "
		}
	}
	reply.Reply = stringa
	return nil

}

func (r RpcMapReduce) EndConnection(arguments EmptyArgument, reply *EmptyReply) error {
	r.DoneWorking <- true
	return nil
}

type EmptyArgument struct{}
type EmptyReply struct{}

type ReduceReply struct {
	Reply string
}

type ReduceMap struct {
	Host     string
	KeyRange [2]int
}
