package fife

type KernelFunction func(args []interface{}, tables map[string]Table)

var kernelInstance  int
var worker          *Worker

func myInstance() int {
    return kernelInstance
}

//RPC calls for non-local data
//used in both worker.go and table.go

type GetArgs {
    Table     int
    Key       string
}

type GetReply {
    Value     interface{}
}

type PutArgs {
    Table     int
    Key       string
    Value     interface{}
}

type PutReply {
    Success   bool
}

//TODO will a flush really be different than a put?
type FlushArgs {
    Table     int
    Key       string
    Value     interface{}
}

type FlushReply {
    Success   bool
}

func myWorker() *Worker {
    return worker
}

const (
    // op types
    CONTAINS = "Contains"
    GET      = "Get"
    PUT      = "Put"
    UPDATE   = "Update"
)

type Op string
