package fife

//Tells us how to treat updates for a table item
type Accumulator struct {
    Init        func(value interface{}) interface{}
    Accumulate  func(originalValue interface{}, newValue interface{}) interface{}
}

var FSumAccumulator = Accumulator {
    Init: func(value interface{}) interface{} { return value },
    Accumulate: func(originalValue interface{}, newValue interface{}) interface{} {
        return originalValue.(float64) + newValue.(float64)
    },
}

//Just returns the first value, for when we don't care which of simultaneous updates
//is returned.
var FirstAccumulator = Accumulator {
    Init: func(value interface{}) interface{} { return value },
    Accumulate: func(originalValue interface{}, newValue interface{}) interface{} {
        return originalValue
    },
}

//Create an accumulator that returns the max of the inputs
//compare returns true if first > second
func CreateMaxAccumulator(compare func(interface{}, interface{}) bool) Accumulator {
  return Accumulator{
    Init: func(value interface{}) interface{} {return value},
    Accumulate: func (initialValue interface{}, newValue interface{}) interface{} {
      if compare(initialValue, newValue) {
        return initialValue
      }
      return newValue
    },
  }
}
