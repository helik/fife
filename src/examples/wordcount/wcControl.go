package wordcount

import (
    "fife"
    "sort"
    "os"
    "strconv"
)

// Control function
func wordCount(f *fife.Fife, numPartitions int, files map[string]string) {
    tables := initTables(numPartitions, nil)

    for k,v := range files {
        tables["documents"].Put(k, v)
    }

    f.Setup(tables)

    numKernels := numPartitions

    f.Run("countWords", numKernels, []interface{}{}, 
        fife.LocalityConstriant{fife.LOCALITY_REQ,"documents"})

    f.Barrier()

    data := f.CollectData("words")

    // get all the keys
    keys := make([]string, 0)
    for k := range data {
        keys = append(keys, k)
    }

    // sort the keys
    sort.Strings(keys)

    // output data to a file in sorted-key order
    file, err := os.Create("results/wc.txt")
    if err != nil { panic(err) }

    for _, k := range keys {
        s := k + ":" + strconv.Itoa(data[k].(int)) + "\n"
        _, err = file.WriteString(s)
        if err != nil { panic(err) }
    }

    err = file.Close()
    if err != nil { panic(err) }
}
