package fife

type KernelFunction func(args []interface{}, tables []Table)

var me int

func myInstance() {
    return me
}