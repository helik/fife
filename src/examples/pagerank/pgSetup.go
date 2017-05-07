package pagerank

import "fife"

func StartWorker(w *fife.Worker, numWorkers int) {
    kernelFunctions := map[string]fife.KernelFunction{"pgKernel":pgKernel}

    tables := initTables(numWorkers*2, w)

    w.Setup(kernelFunctions, tables)
}

func StartFife(f *fife.Fife, numWorkers int) {
    pageRank(f, numWorkers)
}