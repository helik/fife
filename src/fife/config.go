package fife

import (
   "labrpc"
   "sync"
   "testing"
   "runtime"
   crand "crypto/rand"
   "encoding/base64"
   "sync/atomic"
)

func randstring(n int) string {
	b := make([]byte, 2*n)
	crand.Read(b)
	s := base64.URLEncoding.EncodeToString(b)
	return s[0:n]
}

type config struct {
	mu        sync.Mutex
	t         *testing.T
	net       *labrpc.Network
  nWorkers  int
	n         int //total servers = nWorkers + 1
	done      int32 // tell internal threads to die
  Workers   []*Worker //workers
	Fife      *Fife //leader
	applyErr  []string // from apply channel readers
	connected []bool   // whether each server is on the net
	endnames  [][]string    // the port file names each sends to for workers talking to each other
                          //endnames[n-1] is for the fife,  and the fife reference is at the bottom
}

//make a config with n workers and 1 leader fife
func Make_config(t *testing.T, n int) *config {
	runtime.GOMAXPROCS(4)
	cfg := &config{}
	cfg.t = t
	cfg.net = labrpc.MakeNetwork()
    cfg.nWorkers = n
	cfg.n = n + 1
	cfg.applyErr = make([]string, cfg.n)
	cfg.Workers = make([]*Worker, cfg.nWorkers)
  //TODO do we need to init fife here?
	cfg.connected = make([]bool, cfg.n)
	cfg.endnames = make([][]string, cfg.n)

  //TODO this is where we would set unreliable, long delays if we want that

	// create woerkers
	for i := 0; i < cfg.nWorkers; i++ {
		cfg.start_worker(i)
	}
    cfg.start_fife()

	// connect everyone
	for i := 0; i < cfg.n; i++ {
		cfg.connect(i)
	}

	return cfg
}

//
// start a fife
// allocate new outgoing port file names to isolate previous instance of
// this server. since we cannot really kill it.
//
func (cfg *config) start_worker(i int) {

	// a fresh set of outgoing ClientEnd names.
	// so that old crashed instance's ClientEnds can't send.
	cfg.endnames[i] = make([]string, cfg.n)
	for j := 0; j < cfg.n; j++ {
		cfg.endnames[i][j] = randstring(20)
	}

	// a fresh set of ClientEnds.
	ends := make([]*labrpc.ClientEnd, cfg.n)
	for j := 0; j < cfg.n; j++ {
		ends[j] = cfg.net.MakeEnd(cfg.endnames[i][j])
		cfg.net.Connect(cfg.endnames[i][j], j)
	}

	worker := CreateWorker(ends[cfg.n-1], ends[0:cfg.nWorkers], i) //last in ends is fife reference

	cfg.mu.Lock()
	cfg.Workers[i] = worker
	cfg.mu.Unlock()

	svc := labrpc.MakeService(worker)
	srv := labrpc.MakeServer()
	srv.AddService(svc)
	cfg.net.AddServer(i, srv)
}


//
// start or re-start a Raft.
// if one already exists, "kill" it first.
// allocate new outgoing port file names, and a new
// state persister, to isolate previous instance of
// this server. since we cannot really kill it.
//
func (cfg *config) start_fife() {
  //index is last in endnames
  i := cfg.n -1

  // a fresh set of outgoing ClientEnd names.
  // so that old crashed instance's ClientEnds can't send.
  cfg.endnames[i] = make([]string, cfg.n)
  for j := 0; j < cfg.n; j++ {
    cfg.endnames[i][j] = randstring(20)
  }

  // a fresh set of ClientEnds.
  ends := make([]*labrpc.ClientEnd, cfg.n)
  for j := 0; j < cfg.n; j++ {
    ends[j] = cfg.net.MakeEnd(cfg.endnames[i][j])
    cfg.net.Connect(cfg.endnames[i][j], j)
  }

  fi := CreateFife(ends[:cfg.nWorkers]) //last thing in list is reference to ourself

	cfg.mu.Lock()
	cfg.Fife = fi
	cfg.mu.Unlock()

	svc := labrpc.MakeService(fi)
	srv := labrpc.MakeServer()
	srv.AddService(svc)
	cfg.net.AddServer(i, srv)
}

// attach server i to the net.
func (cfg *config) connect(i int) {
	// fmt.Printf("connect(%d)\n", i)

	cfg.connected[i] = true

	// outgoing ClientEnds
	for j := 0; j < cfg.n; j++ {
		if cfg.connected[j] {
			endname := cfg.endnames[i][j]
			cfg.net.Enable(endname, true)
		}
	}

	// incoming ClientEnds
	for j := 0; j < cfg.n; j++ {
		if cfg.connected[j] {
			endname := cfg.endnames[j][i]
			cfg.net.Enable(endname, true)
		}
	}
}

func (cfg *config) cleanup() {
	for i := 0; i < len(cfg.Workers); i++ {
		if cfg.Workers[i] != nil {
			cfg.Workers[i].Kill()
		}
	}
  if cfg.Fife != nil {
    cfg.Fife.Kill()
  }
	atomic.StoreInt32(&cfg.done, 1)
}

//fail the test if not all the master's data is replicated
//on the fife workers, or if the partitions on workers
//don't agree with what we were promised
//Should be run after initial data loaded into master fife
//and workers are configured, but before we run any kernels that would
//change table state
func (cfg *config) CheckDataStore() {
  for tbl_name, table := range(cfg.Fife.tables){
    for partition, worker := range(table.PartitionMap){
      tbl, success := cfg.Workers[worker].tables[tbl_name]
      if !success{
        cfg.t.Fatalf("Worker did not contain expected table")
      }
      _, success = tbl.Store[partition]
      if !success{
        cfg.t.Fatalf("Worker table did not contain expected partition")
      }
      //TODO could check entire data here, but that seems unnecicary 
    }
  }

}
