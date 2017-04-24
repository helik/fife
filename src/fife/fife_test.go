package fife

//

import (
  "testing"
  "fmt"
)


func Test1(t *testing.T){
  workers := 3
  cfg := make_config(t, workers)
  fmt.Println(cfg)
}
