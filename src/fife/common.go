package fife

type KernelFunction func(args []interface{}, tables map[string]Table)

var kernelInstance  int
var worker          *Worker

func myInstance() int {
    return kernelInstance
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