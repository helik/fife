package fife

//Tells us how to treat updates for a table item
type Accumulator struct {
    Init        func(value interface{}) interface{}
    Accumulate  func(originalValue interface{}, newValue interface{}) interface{}
}

func CreateSumAccumulator(sum func(v1 interface{}, v2 interface{}) interface{}) Accumulator {
    return Accumulator {
        Init: func(value interface{}) interface{} { return value },
        Accumulate: func(originalValue interface{}, newValue interface{}) interface{} {
            return sum(originalValue, newValue)
        },
    }
}