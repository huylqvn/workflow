package scheduler

import (
	"context"
	"testing"
	"workflow/cmd/worker"
	"workflow/config"
	"workflow/src/database"
	"workflow/src/domain"

	"github.com/go-co-op/gocron"
	"github.com/huylqvn/httpserver/log"
	"github.com/sirupsen/logrus"
)

func TestScheduler_FetchScheduler(t *testing.T) {
	cfg := config.Load()

	db, err := database.New(cfg)
	if err != nil {
		t.Error(err)
	}
	sqlDb, _ := db.DB()
	defer sqlDb.Close()

	type fields struct {
		cron      *gocron.Scheduler
		log       *logrus.Logger
		taskQueue chan worker.Task
		batchSize int
		db        database.Database
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Scheduler
		wantErr bool
	}{
		// only check query is ok
		{
			name: "TestScheduler_FetchScheduler",
			fields: fields{
				log:       log.Get(),
				taskQueue: make(chan worker.Task, 10),
				batchSize: 10,
				db:        db,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    []domain.Scheduler{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				cron:      tt.fields.cron,
				log:       tt.fields.log,
				taskQueue: tt.fields.taskQueue,
				batchSize: tt.fields.batchSize,
				db:        tt.fields.db,
			}
			_, err := s.FetchScheduler(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scheduler.FetchScheduler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestScheduler_FetchWorkflowTask(t *testing.T) {

	cfg := config.Load()

	db, err := database.New(cfg)
	if err != nil {
		t.Error(err)
	}
	sqlDb, _ := db.DB()
	defer sqlDb.Close()

	type fields struct {
		cron      *gocron.Scheduler
		log       *logrus.Logger
		taskQueue chan worker.Task
		batchSize int
		db        database.Database
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Task
		wantErr bool
	}{
		// only check query is ok
		{
			name: "TestScheduler_FetchWorkflowTask",
			fields: fields{
				log:       log.Get(),
				taskQueue: make(chan worker.Task, 10),
				batchSize: 10,
				db:        db,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    []domain.Task{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				cron:      tt.fields.cron,
				log:       tt.fields.log,
				taskQueue: tt.fields.taskQueue,
				batchSize: tt.fields.batchSize,
				db:        tt.fields.db,
			}
			_, err := s.FetchWorkflowTask(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scheduler.FetchWorkflowTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
