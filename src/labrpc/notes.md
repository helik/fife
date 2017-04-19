RPCs are only capable of passing gob-encoded data, so programs would have to be compiled then passed. This seems like a lot of work when the network is already set up with a master machine and workers - can't we assume that the fife configuration would ensure that the workers knew what kernel code to run? 

Tables and updates can definitely be passed between servers. Just have to define rpcs like the labs did, and make a config file that works for us. 
