package pagerank

import (
    "fife"
    "strings"
    "os"
    "sort"
    "strconv"
)

func pageRank(f *fife.Fife, numPartitions int, input string, numIterations int) {
    propagationFactor := 0.85

    tables := initTables(numPartitions, nil)

    graph := tables["graph"]
    curr := tables["curr"]

    // add initial graph table data & initialize all page ranks to be 1/n 
    //   where n = # of pages = # of entries in graph
    var sites []string
    for _, line := range strings.Split(input, "\n") {
        separated := strings.Split(line, ":")
        siteName := separated[0]

        outlinks := strings.Split(separated[1], ",")

        i := -1
        for j, l := range outlinks {
            if l == siteName { i = j }
        }
        if i >= 0 {
            outlinks = append(outlinks[:i], outlinks[i+1:]...)
        }

        graph.Put(siteName, outlinks)

        sites = append(sites, siteName)
    }

    initialRank := 1/float64(len(sites))
    for _, s := range sites {
        curr.Put(s, initialRank)
    }

    f.Setup(tables)

    for i := 0; i < numIterations; i++ {
        f.Run("pgKernel", numPartitions, []interface{}{propagationFactor}, 
            fife.LocalityConstriant{fife.LOCALITY_REQ,"graph"})

        f.Barrier()

        nextData := f.CollectData("next")

        initTables(numPartitions, nil)

        tables["curr"].AddData(nextData)

        f.Setup(tables)
    }

    data := f.CollectData("curr")

    // get all the keys
    keys := make([]string, 0)
    for k := range data {
        keys = append(keys, k)
    }

    // sort the keys
    sort.Strings(keys)
    // output data to a file in sorted-key order
    file, err := os.OpenFile("results/pg.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
    if err != nil { panic(err) }

    _, err = file.WriteString("PageRank with " + strconv.Itoa(numIterations) + " iterations\n\n")
    if err != nil { panic(err) }

    for _, k := range keys {
        s := k + ":" + strconv.FormatFloat(data[k].(float64),'f',15,64) + "\n"
        _, err = file.WriteString(s)
        if err != nil { panic(err) }
    }

    file.WriteString("\n\n")

    err = file.Close()
    if err != nil { panic(err) }
}