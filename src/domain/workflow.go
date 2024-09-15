package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Status int8

const (
	Inactive  Status = 0
	Active    Status = 1
	Done      Status = 2
	Failed    Status = 3
	Scheduled Status = 4
	Running   Status = 5
)

var weekdayMap = map[string]time.Weekday{
	"sunday":    time.Sunday,
	"monday":    time.Monday,
	"tuesday":   time.Tuesday,
	"wednesday": time.Wednesday,
	"thursday":  time.Thursday,
	"friday":    time.Friday,
	"saturday":  time.Saturday,
}

type Workflow struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`

	Tasks []Task `json:"tasks" gorm:"foreignKey:WorkflowID"`
}

func (Workflow) TableName() string {
	return "workflows"
}

func (w *Workflow) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	w.ID = uuid.New().String()
	w.CreatedAt = &now
	w.UpdatedAt = &now
	return
}

func (w *Workflow) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	w.UpdatedAt = &now
	return
}

type Task struct {
	ID uint64 `json:"id,omitempty" gorm:"primaryKey"`

	WorkflowID   string     `json:"workflow_id"`
	Name         string     `json:"name"`
	Message      string     `json:"message"`
	JobType      string     `json:"job_type"`
	JobTimeValue string     `json:"job_time_value"`
	IsHeadTask   bool       `json:"is_head_task"`
	NextTaskID   uint64     `json:"next_task_id"`
	Status       Status     `json:"status"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

func (t Task) ToRunTime() time.Time {
	// value of period is in number + unit(ms, s, m, h)
	if t.JobType == "period" {
		// ex: 5m to 5*time.Minute
		dur, _ := time.ParseDuration(t.JobTimeValue)
		return time.Now().Add(dur)
	}
	// value of day_of_week is monday, tuesday, wednesday, thursday, friday, saturday, sunday
	if t.JobType == "day_of_week" {
		// find the next day of week
		weekday := time.Now().Weekday()
		target := weekdayMap[t.JobTimeValue]
		dur := int(target - weekday)
		if dur < 0 {
			dur += 7
		}
		return time.Now().Add(time.Duration(dur) * 24 * time.Hour)
	}

	// value of specific_day is 2024-10-01
	if t.JobType == "specific_day" {
		now := time.Now()
		d, _ := time.Parse("2006-01-02 15:04:05", t.JobTimeValue+" "+now.Format("15:04:05"))
		return d
	}

	return time.Now()
}

func (Task) TableName() string {
	return "tasks"
}

func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	t.CreatedAt = &now
	t.UpdatedAt = &now
	return
}

func (t *Task) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	t.UpdatedAt = &now
	return
}

type Scheduler struct {
	ID uint64 `json:"id,omitempty" gorm:"primaryKey"`

	TaskID      uint64     `json:"task_id"`
	RunTime     time.Time  `json:"run_time"`
	PartitionID int        `json:"partition_id"`
	Status      Status     `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func (Scheduler) TableName() string {
	return "schedulers"
}

func (s *Scheduler) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	s.CreatedAt = &now
	s.UpdatedAt = &now
	return
}

func (s *Scheduler) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	s.UpdatedAt = &now
	return
}
