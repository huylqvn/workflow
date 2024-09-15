package worker

import (
	"fmt"
	"workflow/src/database"
	"workflow/src/domain"

	"github.com/huylqvn/httpserver/log"
	"github.com/sirupsen/logrus"
)

type Task interface {
	Run() error
	GetID() uint64
}

type Worker struct {
	id         int
	taskQueue  chan Task
	workerPool chan chan Task
	quitChan   chan bool
	log        *logrus.Logger
	db         database.Database
}

func newWorker(id int, workerPool chan chan Task, db database.Database) Worker {
	return Worker{
		id:         id,
		log:        log.Get(),
		workerPool: workerPool,
		taskQueue:  make(chan Task),
		quitChan:   make(chan bool),
		db:         db,
	}
}

func (w Worker) start() {
	go func() {
		for {
			// Add my taskQueue to the worker pool.
			w.workerPool <- w.taskQueue

			select {
			case job := <-w.taskQueue:
				if err := job.Run(); err != nil {
					w.log.Error("worker error: ", err)
					err = w.db.GetDB().
						Model(&domain.Scheduler{
							ID: job.GetID(),
						}).
						Where("id = ?", job.GetID()).
						Update("status", domain.Failed).Error

					if err != nil {
						w.log.Error("update scheduler error: ", err)
					}
					taskRetry := NewTaskRetry(job)
					go taskRetry.Run()
					// TODO: retry job
					// save to db and retry by scheduler
				} else {
					w.log.Info(fmt.Sprintf("worker %d done job id: %d", w.id, job.GetID()))
				}
			case <-w.quitChan:
				w.log.Info(fmt.Sprintf("worker %d stopping", w.id))
				return
			}
		}
	}()
}

func (w Worker) stop() {
	go func() {
		w.quitChan <- true
	}()
}

func NewDispatcher(taskQueue chan Task, maxWorkers int) *Dispatcher {
	workerPool := make(chan chan Task, maxWorkers)

	return &Dispatcher{
		maxWorkers: maxWorkers,

		workerPool: workerPool,
		taskQueue:  taskQueue,

		log:    log.Get(),
		worker: make([]Worker, 0),
	}
}

type Dispatcher struct {
	maxWorkers int

	workerPool chan chan Task
	taskQueue  chan Task

	worker []Worker
	log    *logrus.Logger
}

func (d *Dispatcher) Run(db database.Database) {
	for i := 0; i < d.maxWorkers; i++ {
		worker := newWorker(i+1, d.workerPool, db)
		worker.start()
		d.worker = append(d.worker, worker)

		d.log.Infof("Worker Started %d", i+1)
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for job := range d.taskQueue {
		go func() {
			workertaskQueue := <-d.workerPool
			workertaskQueue <- job
		}()
	}
}

func (d *Dispatcher) Stop() {
	d.taskQueue = nil

	for _, worker := range d.worker {
		worker.stop()
	}
}
