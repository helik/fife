func initKernelFunctions() []KernelFunction {

}

func startKernel() {
    initTables()

    fife.StartWorker(workers, kernelFunctions, tables)
}