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
