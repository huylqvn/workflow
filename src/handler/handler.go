package handler

import (
	"workflow/config"
	"workflow/src/database"
	"workflow/src/handler/workflow"
	"workflow/src/service"

	"github.com/huylqvn/httpserver"
	"go.elastic.co/apm/module/apmchiv5/v2"
)

func LoadRouter(r httpserver.Router, cfg *config.Config, db database.Database) {
	s := service.NewService(cfg, db)

	// middleware apm tracing
	r.AddMiddleware(apmchiv5.Middleware())

	r.AddPath("/workflow", "POST", workflow.CreateWorkflow(s))
	r.AddPath("/workflow", "GET", workflow.GetListWorkflow(s))

	r.AddPath("/workflow/{id}", "PUT", workflow.UpdateWorkflow(s))
	r.AddPath("/workflow/{id}", "DELETE", workflow.DeleteWorkflow(s))
	r.AddPath("/workflow/{id}", "GET", workflow.GetWorkflow(s))

}
