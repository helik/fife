package fife

import (
  "labrpc"
)

//TODO should tables be an arg to kernel function?
type KernelFunction func(args []interface{}, tables map[string]Table)

var kernelInstance  int //Note: different from worker number 
var worker          *Worker

func myInstance() int {
    return kernelInstance
}

func myWorker() *Worker {
    return worker
}

//The only data a table is ever passed from the fife master
type TableData struct {
  Data                      map[int]map[string]interface{}
  Partitions                map[int]int
}

type RunArgs struct {
    Master                  *labrpc.ClientEnd
    KernelNumber            int
    KernelFunctionName      string
    KernelArgs              []interface{}
}

type RunReply struct {
    Done    bool
}

//Config from fife master to workers
type ConfigArgs struct {
  //data map: string table name to data for that table
  PerTableData       map[string]TableData
}

type ConfigReply struct {
  Success            bool
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
