package fife

import (
  "labrpc"
)

//TODO should tables be an arg to kernel function?
type KernelFunction func(args []interface{}, tables map[string]Table)

var kernelInstance  int
var worker          *Worker

func myInstance() int {
    return kernelInstance
}

func myWorker() *Worker {
    return worker
}


type RunArgs struct {
    Master                  *labrpc.ClientEnd
    KernelNumber            int
    KernelFunctionName      string
    KernelArgs              []interface{}
    //some kind of data thing
}

type RunReply struct {
    Done    bool
}

//Config from fife master to workers
type ConfigArgs struct {
  Data               map[int]map[string]interface{}
  Partitions         map[int]int
}

type ConfigReply struct {

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
