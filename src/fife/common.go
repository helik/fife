package fife

type KernelFunction func(args []interface{}, tables map[string]Table)

var kernelInstance  int
var worker          *Worker

func myInstance() int {
    return kernelInstance
}

//RPC calls for non-local data
//used in both worker.go and table.go

type GetArgs struct {
    Table     int
    Key       string
}

type GetReply struct {
    Value     interface{}
}

type PutArgs struct {
    Table     int
    Key       string
    Value     interface{}
}

type PutReply struct {
    Success   bool
}

//TODO will a flush really be different than a put?
type FlushArgs struct {
    Table     int
    Key       string
    Value     interface{}
}

type FlushReply struct {
    Success   bool
}

func myWorker() *Worker {
    return worker
}

//enum-like listing of possible states
type Op int

const (
	CONTAINS Op = iota
	GET
	PUT
  UPDATE
)
//
// const (
//     // op types
//     CONTAINS = "Contains"
//     GET      = "Get"
//     PUT      = "Put"
//     UPDATE   = "Update"
// )
//
// type Op string
