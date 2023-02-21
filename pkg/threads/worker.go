package threads

import (
	"fmt"
)
//type Worker struct {
//	workerNum int
//	workerCh  chan bool
//}
//
//func CreateWorker(workerNum int) *Worker {
//	return &Worker{
//		workerNum: workerNum,
//		workerCh:  make(chan bool, workerNum),
//	}
//}
//
//// running mult-gorouting
//func (w *Worker) Run(taskFunc func()) {
//	w.workerCh <- true
//	go func() {
//		defer func() { <-w.workerCh }()
//		taskFunc()
//	}()
//}
//
//// Notice: if there is a guard in the main program, then this funcation is not used
//// Used Way: put at the end of the main program，No need to use current funcation if funcation is daemonized.
//func (w *Worker) Wait() {
//	for len(w.workerCh) != 0 && w.workerNum != 0 {
//		log.Infof("There is also gorouting running, the number of runs is %d", len(w.workerCh))
//		time.Sleep(time.Second)
//	}
//	log.Infof("All gorouting is running done")
//}

type Job interface {
	 Do()
}

type JobQueue chan Job

type Worker struct {
	jobChan JobQueue	//每一个worker对象具有JobQueue（队列）属性。
}

type WorkerPool struct {
	WorkerLen int	//协程池的大小
	JobQueue JobQueue	//Job队列，接收外部的数据
	WorkerQueue chan JobQueue //worker队列：处理任务的Go协程队列
}

func NewWorker() Worker {
	return Worker{jobChan: make(chan Job)}
}

func (w Worker) Run(wq chan JobQueue) {
	go func() {
		wq <- w.jobChan
		select {
		case job := <- w.jobChan:
			job.Do()
		}
	}()
}

func NewWorkerPool(workerLen int) *WorkerPool{
	return &WorkerPool{
		WorkerLen:workerLen,
		JobQueue:(make(JobQueue)),
		WorkerQueue: make(chan JobQueue, workerLen),
	}
}

func (wp *WorkerPool) Run() {
	fmt.Println("初始化worker")
	//初始化worker(多个Go程)
	for i := 0; i < wp.WorkerLen; i++ {
		worker := NewWorker()
		worker.Run(wp.WorkerQueue) //开启每一个Go程
	}
	// 循环获取可用的worker,往worker中写job
	go func() {
		for {
			select {
			//将JobQueue中的数据存入WorkerQueue
			case job := <-wp.JobQueue: //线程池中有需要待处理的任务(数据来自于请求的任务) :读取JobQueue中的内容
				worker := <-wp.WorkerQueue //队列中有空闲的Go程   ：读取WorkerQueue中的内容,类型为：JobQueue
				worker <- job              //空闲的Go程执行任务  ：整个job入队列（channel） 类型为：传递的参数（Score结构体）
				//fmt.Println("xxx1:",worker)
				//fmt.Printf("====%T  ;  %T======\n",job,worker,)
			}
		}
	}()
}