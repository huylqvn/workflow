package workflow

import (
	"net/http"
	"workflow/src/domain"
	"workflow/src/service"

	"github.com/go-chi/chi/v5"
	"github.com/huylqvn/httpserver"
	"go.elastic.co/apm/v2"
	"gorm.io/gorm"
)

func CreateWorkflow(s *service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		span, _ := apm.StartSpan(r.Context(), "CreateWorkflow", "endpoint")
		defer span.End()

		req, err := httpserver.ParseBody[WorkflowRequest](r)
		if err != nil {
			apmSendErr(span, err)
			customResponse(w, RequestError, err.Error())
			return
		}

		tx, commit := s.DB.Transaction()
		defer commit(err)
		workflow := &domain.Workflow{
			Name:        req.Name,
			Description: req.Description,
			Status:      domain.Active,
		}
		if err := tx.Create(workflow).Error; err != nil {
			apmSendErr(span, err)
			customResponse(w, InternalServerError, err.Error())
			return
		}

		nextTaskId := uint64(0)
		for i := len(req.Tasks) - 1; i >= 0; i-- {
			isHeadTask := false
			if i == 0 {
				isHeadTask = true
			}
			task := domain.Task{
				WorkflowID:   workflow.ID,
				Name:         req.Tasks[i].Name,
				Message:      req.Tasks[i].Message,
				JobType:      req.Tasks[i].JobType,
				JobTimeValue: req.Tasks[i].JobTimeValue,
				IsHeadTask:   isHeadTask,
				Status:       domain.Active,
				NextTaskID:   nextTaskId,
			}

			if err := tx.Create(&task).Error; err != nil {
				apmSendErr(span, err)
				customResponse(w, InternalServerError, err.Error())
				return
			}
			nextTaskId = task.ID

			workflow.Tasks = append(workflow.Tasks, task)
		}

		customResponse(w, CreateSuccess, workflow)
	}
}

func GetWorkflow(s *service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		span, _ := apm.StartSpan(r.Context(), "GetWorkflow", "endpoint")
		defer span.End()

		id := chi.URLParam(r, "id")
		if id == "" {
			customResponse(w, RequestError, "missing id")
			return
		}

		var workflow domain.Workflow
		if err := s.DB.GetDB().
			Preload("Tasks").
			First(&workflow, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				customResponse(w, NotFound, err.Error())
				return
			}

			apmSendErr(span, err)
			customResponse(w, InternalServerError, err.Error())
			return
		}

		customResponse(w, Success, workflow)
	}
}

func UpdateWorkflow(s *service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		span, _ := apm.StartSpan(r.Context(), "UpdateWorkflow", "endpoint")
		defer span.End()

		req, err := httpserver.ParseBody[WorkflowRequest](r)
		if err != nil {
			apmSendErr(span, err)
			customResponse(w, RequestError, err.Error())
			return
		}

		if req.ID == "" {
			customResponse(w, RequestError, "missing id")
			return
		}

		var workflow domain.Workflow
		if err := s.DB.GetDB().First(&workflow, "id = ?", req.ID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				customResponse(w, NotFound, err.Error())
				return
			}
			apmSendErr(span, err)
			customResponse(w, InternalServerError, err.Error())
			return
		}

		customResponse(w, Success, workflow)
	}
}

func DeleteWorkflow(s *service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		span, _ := apm.StartSpan(r.Context(), "DeleteWorkflow", "endpoint")
		defer span.End()

		var err error

		id := chi.URLParam(r, "id")
		if id == "" {
			customResponse(w, RequestError, "missing id")
			return
		}

		var workflow domain.Workflow
		if err := s.DB.GetDB().First(&workflow, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				customResponse(w, NotFound, err.Error())
				return
			}
			apmSendErr(span, err)
			customResponse(w, InternalServerError, err.Error())
			return
		}

		var idTasks []uint64
		for _, task := range workflow.Tasks {
			idTasks = append(idTasks, task.ID)
		}

		tx, commit := s.DB.Transaction()
		defer commit(err)

		// update status is inactive
		if err := tx.Model(&domain.Task{}).
			Where("id in (?)", idTasks).
			Update("status", domain.Inactive).Error; err != nil {
			apmSendErr(span, err)
			customResponse(w, InternalServerError, err.Error())
			return
		}

		// update workflow is inactive
		if err := tx.Model(&domain.Workflow{}).
			Where("id = ?", id).
			Update("status", domain.Inactive).Error; err != nil {
			apmSendErr(span, err)
			customResponse(w, InternalServerError, err.Error())
			return
		}

		customResponse(w, Success, "delete success")
	}
}

func GetListWorkflow(s *service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		span, _ := apm.StartSpan(r.Context(), "GetListWorkflow", "endpoint")
		defer span.End()

		req, err := httpserver.ParseParam[QueryRequest](r)
		if err != nil {
			apmSendErr(span, err)
			customResponse(w, RequestError, err.Error())
			return
		}
		// set default value
		req.DefaultValue()

		db := s.DB.GetDB().Table("workflows")
		if req.Query != "" {
			db = db.Where("name like ? or description like ?", "%"+req.Query+"%", "%"+req.Query+"%")
		}
		if req.Status != nil {
			db = db.Where("status = ?", *req.Status)
		}
		db = db.Order("id desc")
		db = db.Offset(req.Offset).Limit(req.Limit)

		var workflow []domain.Workflow
		if err := db.
			Preload("Tasks").
			Find(&workflow).Error; err != nil {

			apmSendErr(span, err)
			customResponse(w, InternalServerError, err.Error())
			return
		}

		customResponse(w, Success, workflow)
	}
}
