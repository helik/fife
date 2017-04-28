package fife

//

import (
  "testing"
  "fmt"
)


func TestConfig(t *testing.T){
  workers := 3
  cfg := make_config(t, workers)
  //check workers and fife
  for _, w := range(cfg.workers) {
    if w == nil {
      t.Fatalf("worker not created by config")
    }
  }
  if cfg.fife == nil {
    t.Fatalf("fife not created by config")
  }
  fmt.Println("...passed")
}

/*
Set up some tables, add some data to them, run an rpc call between workers
*/
func TestSetup(t *testing.T){
  cfg := make_config(t, 2) //config with 2 workers
  if len(cfg.workers) != 2 {
    t.Fatalf("unexpected number of workers")
  }
  //shared kernel func between workers
  kernMap := map[string]KernelFunction{"hello":kernel_simple}

  //init workers
  for i, w := range(cfg.workers){
    table := MakeTable(Accumulator{}, Partitioner{}, w, "pasta")
    w.Setup(kernMap, []Table{*table})
    fmt.Println(w)
    //now, add some data
    
  }
}

func TestRPCs(t *testing.T){

}

func kernel_simple(args []interface{}, tables map[string]Table){
  fmt.Println("hello, world")
}
