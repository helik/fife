package pagerank

import "fife"

func pgKernel(kernelInstance int, args []interface{}, tables map[string]*fife.Table) {
    graph := tables["graph"]
    curr := tables["curr"]
    next := tables["next"]

    propagationFactor := args[0].(float64)

    for page, val := range graph.GetPartition(kernelInstance) {
        outlinks := val.([]string)

        rank := curr.Get(page).(float64)
        update := propagationFactor * rank / float64(len(outlinks))

        for _, target := range outlinks {
            next.Update(target, update)
        }
    }

    next.Flush()
}