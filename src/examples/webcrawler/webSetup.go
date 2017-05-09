package webcrawler

import (
  "fife"
)

func StartWorker(w *fife.Worker, numWorkers int) {
    numPartitions := numWorkers //we will have one kernel per worker, and one partition per kernel
    kernelFunctions := map[string]fife.KernelFunction{KERN:fetcherKernel}

    tables := initTables(numPartitions, w)

    w.Setup(kernelFunctions, tables)
}

func StartFife(f *fife.Fife, start_url string, numWorkers int) {
    numPartitions := numWorkers

    webControl(f, numPartitions, start_url)
}
