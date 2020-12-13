package jobx

import (
	"sync"
)

type Worker struct {
	done          *sync.WaitGroup
	readyPool     chan chan Job
	assignedQueue chan Job
	quit          chan bool
}

func NewWorker(readyPool chan chan Job, done *sync.WaitGroup) *Worker {
	return &Worker{
		done:          done,
		readyPool:     readyPool,
		assignedQueue: make(chan Job),
		quit:          make(chan bool),
	}
}

func (w *Worker) Start() {
	w.done.Add(1)
	go func(worker *Worker) {
		for {
			worker.readyPool <- worker.assignedQueue
			select {
			case job := <-worker.assignedQueue:
				job.Process()
			case <-worker.quit:
				worker.done.Done()
				return
			}
		}
	}(w)
}

func (w *Worker) Stop() {
	w.quit <- true
}
