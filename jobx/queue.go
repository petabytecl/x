package jobx

import (
	"errors"
	"sync"
)

type Queue struct {
	internal          chan Job
	readyPool         chan chan Job
	workers           []*Worker
	stoppedMutex      *sync.Mutex
	stoppedWorkers    *sync.WaitGroup
	stoppedDispatcher *sync.WaitGroup
	stopped           bool
	quit              chan bool
}

func NewQueue(maxWorkers int) *Queue {
	stoppedWorkers := &sync.WaitGroup{}
	readyPool := make(chan chan Job, maxWorkers)
	workers := make([]*Worker, maxWorkers)

	for i := 0; i < maxWorkers; i++ {
		workers[i] = NewWorker(readyPool, stoppedWorkers)
	}

	return &Queue{
		internal:          make(chan Job),
		readyPool:         readyPool,
		workers:           workers,
		stoppedDispatcher: &sync.WaitGroup{},
		stoppedWorkers:    stoppedWorkers,
		quit:              make(chan bool),
	}
}

func (q *Queue) Start() {
	for i := 0; i < len(q.workers); i++ {
		q.workers[i].Start()
	}

	go q.dispatch()

	q.stoppedMutex = &sync.Mutex{}
	q.setStopped(false)
}

func (q *Queue) Stop() {
	q.setStopped(true)

	q.quit <- true
	q.stoppedDispatcher.Wait()

	close(q.internal)
}

func (q *Queue) Stopped() bool {
	q.stoppedMutex.Lock()
	lock := q.stopped
	q.stoppedMutex.Unlock()

	return lock
}

func (q *Queue) setStopped(lock bool) {
	q.stoppedMutex.Lock()
	q.stopped = lock
	q.stoppedMutex.Unlock()
}

func (q *Queue) dispatch() {
	q.stoppedDispatcher.Add(1)

	for {
		select {
		case job := <-q.internal:
			workerChan := <-q.readyPool
			workerChan <- job
		case <-q.quit:
			for i := 0; i < len(q.workers); i++ {
				q.workers[i].Stop()
			}
			q.stoppedWorkers.Wait()
			q.stoppedDispatcher.Done()
			return
		}
	}
}

func (q *Queue) Submit(job Job) (bool, error) {
	q.stoppedMutex.Lock()
	if q.stopped {
		q.stoppedMutex.Unlock()
		return false, errors.New("queue is stopped")
	}

	q.internal <- job

	q.stoppedMutex.Unlock()

	return true, nil
}
