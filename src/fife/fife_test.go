package fife

//

import (
  "testing"
  "fmt"
)


func Test1(t *testing.T){
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

func TestStartWorkers(t *testing.T){

}
