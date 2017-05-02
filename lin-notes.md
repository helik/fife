
fife.CreateTable and MakeTable in Table have very similar arguments. Perhaps fife.CreateTable
should take in a bare-bones table as created by MakeTable?  

idea for extension - could we double up which workers have which partitions, so in the case of a
worker crashing or a network partition, other workers don't have to wait for the master to
re-partition that data and re-configure workers? 
