package wordcount

import (
    "fife"
    "unicode"
    "strings"
    "fmt"
)

// Kernel function
func countWords(args []interface{}, tables map[string]*fife.Table) {
    documents := tables["documents"]
    words := tables["words"]
    isNotALetter := func(c rune) bool {
        return !unicode.IsLetter(c)
    }

    // look at all documents in this partition
    for _, doc := range documents.GetPartition(fife.MyInstance()) {
        fmt.Println(doc)
        // for all the words in the document
        for _, word := range strings.FieldsFunc(getDocValue(doc), isNotALetter) {
            // increment the number of words in store
            words.Update(word, 1)
        }
    }

    words.Flush()
}