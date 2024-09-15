package scheduler

import (
	"context"
	"fmt"
	"time"
	"workflow/cmd/worker"
	"workflow/src/database"
	"workflow/src/domain"

	"github.com/go-co-op/gocron"
	"github.com/huylqvn/httpserver/log"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/v2"
	"gorm.io/gorm/clause"
)

type Scheduler struct {
	cron      *gocron.Scheduler
	log       *logrus.Logger
	taskQueue chan worker.Task
	batchSize int
	db        database.Database
}

func NewScheduler(taskQueue chan worker.Task, batchSite int, db database.Database) *Scheduler {
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	cron := gocron.NewScheduler(loc)
	return &Scheduler{
		cron:      cron,
		log:       log.Get(),
		taskQueue: taskQueue,
		batchSize: batchSite,
		db:        db,
	}
}

func (s *Scheduler) TimeDelay() *Scheduler {
	f := func() {
		tx := apm.DefaultTracer().StartTransaction("Scheduler TimeDelay", "request")
		defer tx.End()
		ctx := apm.ContextWithTransaction(context.Background(), tx)

		tasks, err := s.FetchWorkflowTask(ctx)
		if err != nil {
			tx.Result = "error"
			apm.CaptureError(ctx, err)
			s.log.Error("fetch task error", err)
			return
		}

		if len(tasks) == 0 {
			// don't send to apm
			tx.Discard()
			return
		}
		tx.Context.SetLabel("num_of_task", len(tasks))

		for _, task := range tasks {
			s.log.Info(fmt.Sprintf("create scheduler for task id: %d, name: %s", task.ID, task.Name))
			tx.Context.SetLabel(fmt.Sprintf("create_scheduler_%d", task.ID), task.Name)

			// create scheduler
			ch := domain.Scheduler{
				TaskID:      task.ID,
				RunTime:     task.ToRunTime(),
				Status:      domain.Active,
				PartitionID: 0,
			}

			if err := s.CreateScheduler(ch); err != nil {
				tx.Result = "error"
				apm.CaptureError(ctx, err)
				s.log.Error("create scheduler error", err)
				return
			}
		}
		tx.Result = "success"

	}
	s.cron.Every(5).Seconds().Do(f)
	return s
}

func (s *Scheduler) CreateScheduler(ch domain.Scheduler) error {
	return s.db.GetDB().Create(&ch).Error
}

func (s *Scheduler) PushScheduler() *Scheduler {
	f := func() {
		tx := apm.DefaultTracer().StartTransaction("Scheduler PushScheduler", "request")
		defer tx.End()
		ctx := apm.ContextWithTransaction(context.Background(), tx)

		schedulers, err := s.FetchScheduler(ctx)
		if err != nil {
			tx.Result = "error"
			apm.CaptureError(ctx, err)
			s.log.Error("fetch scheduler error", err)
			return
		}

		if len(schedulers) == 0 {
			// don't send to apm
			tx.Discard()
			return
		}

		tx.Context.SetLabel("num_of_scheduler", len(schedulers))

		for _, ch := range schedulers {
			// alter push to message queue for scale solution
			t := worker.NewTaskSendNoti(ctx, ch.ID, s.db)
			s.taskQueue <- t
		}

		tx.Result = "success"
	}
	s.cron.Every(1).Seconds().Do(f)
	return s
}

func (s *Scheduler) FetchWorkflowTask(ctx context.Context) ([]domain.Task, error) {
	span, _ := apm.StartSpan(ctx, "FetchWorkflowTask", "func")
	defer span.End()

	var tasks []domain.Task

	if err := s.db.GetDB().
		Model(&tasks).
		Clauses(clause.Returning{}).
		Where("status = ?", domain.Active).
		Where("is_head_task = ?", true).
		Update("status", domain.Scheduled).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Scheduler) FetchScheduler(ctx context.Context) ([]domain.Scheduler, error) {
	span, _ := apm.StartSpan(ctx, "FetchScheduler", "func")
	defer span.End()

	var schedulers []domain.Scheduler
	if err := s.db.GetDB().Model(&schedulers).
		Clauses(clause.Returning{}).
		Where("status = ?", domain.Active).
		Where("run_time < ?", time.Now()).
		Order("created_at desc").
		Update("status", domain.Running).Error; err != nil {
		return nil, err
	}
	return schedulers, nil
}

func (s *Scheduler) Start() {
	s.cron.SetMaxConcurrentJobs(1, gocron.WaitMode)
	s.cron.StartAsync()
	s.log.Info("start scheduler")
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	s.log.Info("stop scheduler")
}
