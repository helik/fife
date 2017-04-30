Test 1 -- check basic kernel function
   -kernel function just prints its kernel instance number
   -tests: machine setup, scheduling, run RPC, myInstance()

Test 2 -- check table partition & setup
   -master sets up a table with data already in it & sends it to workers
   -the kernel function will just print the table data it received
   -tests: master table setup, worker table setup, config RPC

Test 3 -- check local gets
   -same setup as test 2, kernel function executes and prints some local gets

Test 4 -- check local puts
   -same setup as test 2, kernel function  executes some local puts

Test 5 -- check local updates
   -same setup as test 2, kernel function executes some local updates

Test 6 -- check remote gets

Test 7 -- check remote puts

Test 8 -- check remote updates & flush

Test 9 -- word count example

Test 10 -- page rank example

Test 11 -- web crawler example

Test 12 -- task stealing (probably easiest with wc example)
   -slow down one of the workers, steal a task from the worker and send to another

Test 13 -- task stealing & partition migration
   -same as test 12, but migrate partition when stealing a task

Test 14 -- retest word count

Test 15 -- retest page rank

Test 16 -- retest web crawler

Test 17+ snapshotting...
