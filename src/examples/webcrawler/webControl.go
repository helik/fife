package webcrawler

import (
  "fife"
  "encoding/gob"
)

// Control function
func webControl(f *fife.Fife, numPartitions int, start_url string) {
    numKernels := numPartitions
    tables := initTables(numPartitions, nil)

    gob.Register(TOFETCH)  //Need to register all interface{} types that we'll send over rpc

    tables[URL_TABLE].Put(start_url, TOFETCH) //seed

    f.Setup(tables)

    f.Run(KERN, numKernels, []interface{}{},
        fife.LocalityConstriant{fife.LOCALITY_REQ, URL_TABLE})

    f.Barrier()
}
 //Need to:
/*
 gob.Register any interface type we want to send over rpc.


 */
