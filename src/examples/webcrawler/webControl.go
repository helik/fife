package webcrawler

import (
  "fife"
)

// Control function
func webControl(f *fife.Fife, numPartitions int) {
    numKernels := numPartitions

    f.Run(KERN, numKernels, []interface{}{},
        fife.LocalityConstriant{fife.LOCALITY_REQ, URL_TABLE})

    f.Barrier()
}
