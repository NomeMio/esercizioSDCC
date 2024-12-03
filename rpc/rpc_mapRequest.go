package rpc

import (
	"encoding/json"
	"esercizioSDCC/utilis"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// struttura su cui si registra del rpc
type MapRequest struct {
	frequency map[int][]int
	lock      sync.Locker
	address   string
}

// arogmento della chiamata MapGetResult che rappresenta lo shard del file
type MapArguemnt struct {
	InputString string
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
	return newObj
}

func (r MapRequest) MapAppending(key int, value int) {
	r.lock.Lock()
	r.frequency[key] = append(r.frequency[key], value)
	r.lock.Unlock()
}

// effettua la mappatura dello shard nella mappa frequency e ritorna il valore massim e minimo fra le chiavi della mappa
func (r MapRequest) MapGetResult(arguemnt MapArguemnt, reply *MapReply) error {
	//TODO vedere cosa fare in caso di chiamate duplicate
	if len(r.frequency) != 0 {
		r.frequency = make(map[int][]int)
	}
	tokens := strings.Fields(arguemnt.InputString)
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

func (r MapRequest) ShowMyRange(arguemnt []ReduceMap, reply *ReduceReply) error {
	for _, item := range arguemnt {
		if item.Host == r.address {
			fmt.Printf("my range is %v\n", item.KeyRange)
		}
	}
	return nil
}

type ReduceReply struct {
	Reply string
}

type ReduceMap struct {
	Host     string
	KeyRange [2]int
}
