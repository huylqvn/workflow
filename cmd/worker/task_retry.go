package worker

import (
	"time"

	"github.com/huylqvn/httpserver/log"
)

type TaskRetry struct {
	Task       Task
	RetryCount int
	Delay      int
}

func NewTaskRetry(task Task) *TaskRetry {
	return &TaskRetry{
		Task:  task,
		Delay: 10,
	}
}

func (t *TaskRetry) Run() error {
	if t.RetryCount >= 3 {
		return nil
	}
	time.Sleep(time.Duration(t.Delay) * time.Second)
	t.RetryCount++
	log.Get().Info("task retry: ", t.RetryCount, " scheduler id: ", t.Task.GetID())
	return t.Task.Run()
}

func (t *TaskRetry) GetID() uint64 {
	return t.Task.GetID()
}
