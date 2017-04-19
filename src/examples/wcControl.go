func initTables() []Table {
    // call fife to init tables
}

func startControl() {
    initTables()

    fife.StartControl(workers, tables)
}