package worker

import (
	"context"
	"fmt"
	"workflow/config"
	"workflow/src/database"
	"workflow/src/domain"
	"workflow/src/service/bot"

	"github.com/huylqvn/httpserver/log"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TaskSendNoti struct {
	TaskID uint64
	Ctx    context.Context
	DB     database.Database

	log *logrus.Logger
}

func NewTaskSendNoti(ctx context.Context, taskId uint64, db database.Database) *TaskSendNoti {
	return &TaskSendNoti{
		TaskID: taskId,
		Ctx:    ctx,
		DB:     db,

		log: log.Get(),
	}
}

func (t *TaskSendNoti) GetID() uint64 {
	return t.TaskID
}

func (t *TaskSendNoti) Run() error {
	t.log.Infof("Executing task scheduler: %d", t.TaskID)

	span, _ := apm.StartSpan(t.Ctx, "TaskSendNoti", "func")
	defer span.End()
	span.Context.SetLabel("scheduler_id", t.TaskID)

	var err error
	tx, commit := t.DB.Transaction()
	defer commit(err)

	// get scheduler
	ch, err := t.GetScheduler(tx)
	if err != nil {
		return err
	}

	// get task
	task, err := t.GetTask(tx, ch.TaskID)
	if err != nil {
		return err
	}

	span.Context.SetLabel("task_name", task.Name)
	t.log.Info(fmt.Sprintf("task name: %s", task.Name))

	// send message to telegram, rate limit 30 per minute

	bot := bot.Get()
	cfg := config.Load()
	err = bot.SendGroup(cfg.GroupID, task.Message)
	if err != nil {
		apm.CaptureError(t.Ctx, err)
		t.log.Error("send message error", err)
		return err
	}

	// save scheduler status
	err = t.UpdateSchedulerStatus(tx)
	if err != nil {
		return err
	}

	// create next task
	return t.CreateNextTask(tx, task.NextTaskID)
}

func (t *TaskSendNoti) GetScheduler(tx *gorm.DB) (*domain.Scheduler, error) {
	var ch domain.Scheduler
	if err := tx.First(&ch, "id = ?", t.TaskID).Error; err != nil {
		apm.CaptureError(t.Ctx, err)
		t.log.Error("fetch scheduler error", err)
		return nil, err
	}

	return &ch, nil
}

func (t *TaskSendNoti) GetTask(tx *gorm.DB, id uint64) (*domain.Task, error) {
	var task domain.Task
	if err := tx.Model(&task).
		First(&task, "id = ?", id).Error; err != nil {
		apm.CaptureError(t.Ctx, err)
		t.log.Error("fetch task error", err)
		return nil, err
	}

	return &task, nil
}

func (t *TaskSendNoti) UpdateSchedulerStatus(tx *gorm.DB) error {
	if err := tx.
		Model(&domain.Scheduler{}).
		Where("id = ?", t.TaskID).
		Update("status", domain.Done).Error; err != nil {
		apm.CaptureError(t.Ctx, err)
		t.log.Error("update scheduler error", err)
		return err
	}

	return nil
}

func (t *TaskSendNoti) CreateNextTask(tx *gorm.DB, nextTextId uint64) error {
	if nextTextId != 0 {
		var nextTask domain.Task
		if err := tx.Model(&nextTask).
			Clauses(clause.Returning{}).
			Where("id = ?", nextTextId).
			Where("status = ?", domain.Active).
			Update("status", domain.Scheduled).Error; err != nil {
			apm.CaptureError(t.Ctx, err)
			t.log.Error("update task error", err)
			return err
		}

		nextCh := domain.Scheduler{
			TaskID:      nextTask.ID,
			RunTime:     nextTask.ToRunTime(),
			Status:      domain.Active,
			PartitionID: 0,
		}

		if err := tx.Create(&nextCh).Error; err != nil {
			apm.CaptureError(t.Ctx, err)
			t.log.Error("create scheduler error", err)
			return err
		}
		t.log.Info("create scheduler next task: ", nextTask.Name)
	}

	return nil
}
