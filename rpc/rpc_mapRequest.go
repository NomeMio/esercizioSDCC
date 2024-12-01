package rpc

import (
	"encoding/json"
	"esercizioSDCC/utilis"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type MapRequest struct {
	frequency map[int][]int
	lock      sync.Locker
	address   string
}
type MapArguemnt struct {
	InputString string
}
type MapReply struct {
	MaxValue int
	MinValue int
}

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

func (r MapRequest) MapGetResult(arguemnt MapArguemnt, reply *MapReply) error {
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
