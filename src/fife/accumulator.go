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