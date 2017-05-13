package fife

//

import (
  "testing"
  "fmt"
  "time"
  "strconv"
)

const name string = "first_table"

//data to initialize fife table with
var data map[string]interface{} = map[string]interface{}{"apple":1, "banana":100, "zebra":200, "cat":40,"annie":2.2}

var data_large map[string]interface{} = map[string]interface{}{"0":0, "1":1, "2":2, "3":3,
              "4":4, "5":5, "6":6, "7":7, "8":8, "9":9, "10":10, "11":11, "12":12}

func TestConfig(t *testing.T){
  workers := 3
  cfg := Make_config(t, workers)
  //check workers and fife
  for _, w := range(cfg.Workers) {
    if w == nil {
      t.Fatalf("worker not created by config")
    }
  }
  if cfg.Fife == nil {
    t.Fatalf("fife not created by config")
  }
  fmt.Println("...passed")
}

/*
Set up workers, run an rpc call from fife to workers, call simple kernel on workers.
Workers should each print hello, world
Note: worker run called by test, not rpc from fife
*/
func TestSetup(t *testing.T){
  cfg := Make_config(t, 2) //config with 2 workers
  if len(cfg.Workers) != 2 {
    t.Fatalf("unexpected number of workers")
  }

  //shared kernel func between workers
  kernName := "hello"
  kernMap := map[string]KernelFunction{kernName:kernel_simple}

  //init workers
  for _, w := range(cfg.Workers){
    table := MakeTable(name, Accumulator{}, Partitioner{}, 0,w) //not using accumulator or partitioner for this test
    w.Setup(kernMap, map[string]*Table{name:table})
  }

  //call an rpc from master
  ok := cfg.Fife.configWorkers()
  if !ok{
    t.Fatalf("Some config rpcs failed")
  }

  //run kernel functions
  for i, w := range(cfg.Workers){
    args := &RunArgs{}
    reply := &Reply{}
    //args.Master = w.fife
    args.KernelNumber = i
    args.KernelFunctionName = kernName
    args.KernelArgs = []interface{}{i}
    w.Run(args, reply)
  }
  fmt.Println("...passed")
}

func kernel_simple(kernelInstance int, args []interface{}, tables map[string]*Table){
  fmt.Printf("hello, world. Kernel instance %v\n", kernelInstance)
}

func kernel_locality(kernelInstance int, args []interface{}, tables map[string]*Table){
  fmt.Printf("kernel %v. table store %v\n",kernelInstance, tables[name].Store)
}

//kernel where one kernel instance waits
func kernel_wait(kernelInstance int, args []interface{}, tables map[string]*Table){
  if kernelInstance == 0 {
    time.Sleep(time.Second)
  }
}

//partitions into however many int imputs data_large has
func partition_many(key string) int{
  val, _ := strconv.Atoi(key)
  return val
}

func partition_simple(key string) int{
  switch key[0]{
  case 'a':
    return 0
  case 'b':
    return 1
  case 'c':
    return 2
  default:
    return 3 //n partitions = 4
  }
}


//Create a table in master fife, check that it partitioned correctly.
//Assign partitions to workers using partitionTables
//Then, send those tables to workers using configWorkers.
//Note that these last two steps are performed together during Run, just not testing run quite yet
func TestFifeTable(t *testing.T){
  cfg := Make_config(t, 2) //config with 2 workers

  tableName := "table1"

  //init workers
  for _, w := range(cfg.Workers){
    table := MakeTable(tableName, Accumulator{}, Partitioner{partition_simple}, 4, w) //not using accumulator or partitioner for this test
    w.Setup(make(map[string]KernelFunction), map[string]*Table{tableName:table})
  }

  table := MakeTable(tableName, Accumulator{}, Partitioner{partition_simple}, 4, nil)
  table.AddData(data)

  cfg.Fife.Setup(map[string]*Table{tableName:table})

  if len(cfg.Fife.tables[tableName].Store) != 4 {
    t.Fatalf("Expected 4 partitions from simple_partition")
  }
  if len(cfg.Fife.tables[tableName].Store[0]) != 2 {
    t.Fatalf("Expected two \"a\" names to be partitioned together")
  }

  cfg.Fife.partitionTables()
  if len(cfg.Fife.tables[tableName].PartitionMap) != 4{
    t.Fatalf("Expected 4 mapped partitions")
  }
  for partition, worker := range(cfg.Fife.tables[tableName].PartitionMap){
    if partition > 3 || worker > 1 {
      t.Fatalf("Expected largest partition 3 and largest worker 1")
    }
  }

  ok := cfg.Fife.configWorkers()
  if !ok{
    t.Fatalf("Some config rpcs failed")
  }

  for _, w := range(cfg.Workers){
    if len(w.tables[tableName].PartitionMap) != 4{
      t.Fatalf("Expected partition map to transfer from master to workers")
    }
    fmt.Println(w.tables[tableName].Store)
  }

  fmt.Println("...passed")
}

func TestFifeRun(t *testing.T){
  cfg := Make_config(t, 3) //config with 3 workers

  tableName := "table1"

  //init workers
  kernName := "hello" //shared kernel func between workers
  kernMap := map[string]KernelFunction{kernName:kernel_simple}
  for _, w := range(cfg.Workers){
    table := MakeTable(tableName, Accumulator{}, Partitioner{partition_simple}, 4, w) //not using accumulator or partitioner for this test
    w.Setup(kernMap, map[string]*Table{tableName:table})
  }

  table := MakeTable(tableName, Accumulator{}, Partitioner{partition_simple}, 4, nil)
  table.AddData(data)

  cfg.Fife.Setup(map[string]*Table{tableName:table})

  cfg.Fife.Run("hello", 6, []interface{}{}, LocalityConstriant{NONE, ""})

  cfg.Fife.Barrier()

  cfg.CheckDataStore()
}

func TestLocality(t *testing.T){
  cfg := Make_config(t, 4) //config with same # workers as # partitions

  tableName := name

  //init workers
  kernName := "locality" //shared kernel func between workers
  kernMap := map[string]KernelFunction{kernName:kernel_locality}
  for _, w := range(cfg.Workers){
    table := MakeTable(tableName, Accumulator{}, Partitioner{partition_simple}, 4, w) //not using accumulator or partitioner for this test
    w.Setup(kernMap, map[string]*Table{tableName:table})
  }

  table := MakeTable(tableName, Accumulator{}, Partitioner{partition_simple}, 4, nil)
  table.AddData(data)

  cfg.Fife.Setup(map[string]*Table{tableName:table})

  fmt.Printf("Should print: \nKernel 0, a keys; kernel 1; b keys; kernel 2; c keys, kernel 3, others\n")
  cfg.Fife.Run("locality", 4, []interface{}{}, LocalityConstriant{LOCALITY_REQ, tableName})

  cfg.Fife.Barrier()

  cfg.CheckDataStore()
}

func TestRePartition(t *testing.T){
  cfg := Make_config(t, 3) //config with same # workers as # partitions

  tableName := name

  //init workers
  kernName := "time" //shared kernel func between workers
  kernMap := map[string]KernelFunction{kernName:kernel_locality}
  for _, w := range(cfg.Workers){
    table := MakeTable(tableName, Accumulator{}, Partitioner{partition_many}, 13, w) //not using accumulator or partitioner for this test
    w.Setup(kernMap, map[string]*Table{tableName:table})
  }

  table := MakeTable(tableName, Accumulator{}, Partitioner{partition_many}, 13, nil)
  table.AddData(data_large)

  cfg.Fife.Setup(map[string]*Table{tableName:table})

  cfg.CheckDataStore()

  fmt.Println("Should repartition")
  cfg.Fife.Run("time", 13, []interface{}{}, LocalityConstriant{LOCALITY_REQ, tableName})

  cfg.Fife.Barrier()

  cfg.CheckDataStore()

}
