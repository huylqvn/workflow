package workflow

import (
	"errors"
	"workflow/src/domain"
)

type WorkflowRequest struct {
	ID          string        `json:"id,omitempty"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Tasks       []TaskRequest `json:"tasks"`
}

type TaskRequest struct {
	ID           uint64 `json:"id,omitempty"`
	Name         string `json:"name"`
	Message      string `json:"message"`
	JobType      string `json:"job_type"`
	JobTimeValue string `json:"job_time_value"`
}

func (r WorkflowRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	for _, t := range r.Tasks {
		if t.JobType == "" || t.JobTimeValue == "" {
			return errors.New("job_type and job_time_value is required")
		}
		if t.Message == "" {
			return errors.New("message is required")
		}
	}
	return nil
}

type QueryRequest struct {
	Query  string         `json:"query"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
	Status *domain.Status `json:"status"`
}

func (r QueryRequest) Validate() error {
	return nil
}

func (r *QueryRequest) DefaultValue() {
	if r.Limit == 0 {
		r.Limit = 10
	}

	if r.Offset == 0 {
		r.Offset = 0
	}
}
