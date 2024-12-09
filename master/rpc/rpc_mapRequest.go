package rpc

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

type EmptyArgument struct{}
type EmptyReply struct{}

type ReduceReply struct {
	Reply string
}

type ReduceMap struct {
	Host     string
	KeyRange [2]int
}
