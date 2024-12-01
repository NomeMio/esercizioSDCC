package rpc

import (
	"encoding/json"
	"esercizioSDCC/utilis"
	"strings"
)

type MapRequest struct {
	frequency map[string]int
}
type MapArguemnt struct {
	InputString string
}
type MapReply struct {
	MaxValue int
	MinValue int
}

func (r MapRequest) MapGetResult(arguemnt MapArguemnt, reply *MapReply) error {
	r.frequency = make(map[string]int)
	tokens := strings.Fields(arguemnt.InputString)
	for _, token := range tokens {
		r.frequency[token] += 1
	}
	jsonString, err := json.Marshal(r.frequency)
	utilis.CheckError(err)
	print(string(jsonString))
	return nil
}
