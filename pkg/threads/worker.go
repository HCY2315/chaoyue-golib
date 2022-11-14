package threads

import (
	"time"

	"github.com/HCY2315/chaoyue-golib/log"
)

type worker struct {
	workerNum int
	workerCh  chan bool
}

func CreateWorker(workerNum int) *worker {
	return &worker{
		workerNum: workerNum,
		workerCh:  make(chan bool, workerNum),
	}
}

func (w *worker) Run(taskFunc func()) {
	w.workerCh <- true
	go func() {
		defer func() { <-w.workerCh }()
		taskFunc()
	}()
}

// Notice: if there is a guard in the main program, then this funcation is not used
// Used Way: put at the end of the main program
func (w *worker) Wait() {
	for len(w.workerCh) != 0 && w.workerNum != 0 {
		time.Sleep(time.Second * 1)
		log.Infof("There is also gorouting running, the number of runs is %d", len(w.workerCh))
	}
	log.Infof("All gorouting is running done")
}
