package experiments

import (
  "labrpc"
  "fmt"
)

func torun() {
  network := labrpc.MakeNetwork()
  fmt.Println(network)
}
