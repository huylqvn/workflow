package worker

import (
	"context"
	"reflect"
	"testing"
	"workflow/config"
	"workflow/src/database"
	"workflow/src/domain"

	"github.com/huylqvn/httpserver/log"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func TestTaskSendNoti_Run(t *testing.T) {

	cfg := config.Load()

	db, err := database.New(cfg)
	if err != nil {
		t.Error(err)
	}
	sqlDb, _ := db.DB()
	defer sqlDb.Close()

	type fields struct {
		TaskID uint64
		Ctx    context.Context
		DB     database.Database
		log    *logrus.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// happy case
		{
			name: "happy case send hello",
			fields: fields{
				TaskID: 1,
				Ctx:    context.Background(),
				DB:     db,
				log:    log.Get(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TaskSendNoti{
				TaskID: tt.fields.TaskID,
				Ctx:    tt.fields.Ctx,
				DB:     tt.fields.DB,
				log:    tt.fields.log,
			}
			if err := tr.Run(); (err != nil) != tt.wantErr {
				t.Errorf("TaskSendNoti.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaskSendNoti_GetScheduler(t *testing.T) {

	cfg := config.Load()

	db, err := database.New(cfg)
	if err != nil {
		t.Error(err)
	}
	sqlDb, _ := db.DB()
	defer sqlDb.Close()

	type fields struct {
		TaskID uint64
		Ctx    context.Context
		DB     database.Database
		log    *logrus.Logger
	}
	type args struct {
		tx *gorm.DB
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Scheduler
		wantErr bool
	}{
		// get scheduler id = 1
		{
			name: "get scheduler",
			fields: fields{
				TaskID: 1,
				Ctx:    context.Background(),
				DB:     db,
				log:    log.Get(),
			},
			args: args{
				tx: db.GetDB(),
			},
			want: &domain.Scheduler{
				ID: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TaskSendNoti{
				TaskID: tt.fields.TaskID,
				Ctx:    tt.fields.Ctx,
				DB:     tt.fields.DB,
				log:    tt.fields.log,
			}
			got, err := tr.GetScheduler(tt.args.tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskSendNoti.GetScheduler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.ID, tt.want.ID) {
				t.Errorf("TaskSendNoti.GetScheduler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskSendNoti_GetTask(t *testing.T) {

	cfg := config.Load()

	db, err := database.New(cfg)
	if err != nil {
		t.Error(err)
	}
	sqlDb, _ := db.DB()
	defer sqlDb.Close()

	type fields struct {
		TaskID uint64
		Ctx    context.Context
		DB     database.Database
		log    *logrus.Logger
	}
	type args struct {
		tx *gorm.DB
		id uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Task
		wantErr bool
	}{
		// get task id = 1
		{
			name: "get task",
			fields: fields{
				TaskID: 1,
				Ctx:    context.Background(),
				DB:     db,
				log:    log.Get(),
			},
			args: args{
				tx: db.GetDB(),
				id: 1,
			},
			want: &domain.Task{
				ID: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TaskSendNoti{
				TaskID: tt.fields.TaskID,
				Ctx:    tt.fields.Ctx,
				DB:     tt.fields.DB,
				log:    tt.fields.log,
			}
			got, err := tr.GetTask(tt.args.tx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskSendNoti.GetTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.ID, tt.want.ID) {
				t.Errorf("TaskSendNoti.GetTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskSendNoti_UpdateSchedulerStatus(t *testing.T) {

	cfg := config.Load()

	db, err := database.New(cfg)
	if err != nil {
		t.Error(err)
	}
	sqlDb, _ := db.DB()
	defer sqlDb.Close()

	type fields struct {
		TaskID uint64
		Ctx    context.Context
		DB     database.Database
		log    *logrus.Logger
	}
	type args struct {
		tx *gorm.DB
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// update scheduler id = 1
		{
			name: "update scheduler",
			fields: fields{
				TaskID: 1,
				Ctx:    context.Background(),
				DB:     db,
				log:    log.Get(),
			},
			args: args{
				tx: db.GetDB(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TaskSendNoti{
				TaskID: tt.fields.TaskID,
				Ctx:    tt.fields.Ctx,
				DB:     tt.fields.DB,
				log:    tt.fields.log,
			}
			if err := tr.UpdateSchedulerStatus(tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("TaskSendNoti.UpdateSchedulerStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
