import "fife"
import "fife/table"
import "hash/fnv"

type Value struct {
    value   int
}

func makeValue(val int) Value {
    return Value{val}
}

func getValue(val interface{}) int {
    return val.(Value).value
}

func partitioner(key string) int {
    return 1
}

func countWords(args []interface{}) {
    documents := args[0].(Table)
    words := args[1].(Table)
    isNotALetter := func(c rune) bool {
        return !unicode.IsLetter(c)
    }
    // look at all documents in this partition
    for _, doc := range documents.getPartition(myInstance()) {
        // for all the words in the document
        for _, word := range strings.FieldsFunc(doc, isNotALetter) {
            // increment the number of words in store
            words.update(word, 1)
        }
    }
}

func wordCount() {
    ff := fife.Make()

    documents := ff.CreateTable(1, Accumulator{},
        Partitioner{partitioner})

    words := ff.CreateTable(1, 
        Accumulator{
            func(value interface{}) interface{} {return value},
            func(original interface{}, newVal interface{}) interface{} {
                return makeValue(getValue(original) + getValue(newVal))
                },
            },
        Partitioner{partitioner})

    ff.Run(countWords, []interface{}{documents, words})

    ff.Barrier()

    fmt.Println(words)
}