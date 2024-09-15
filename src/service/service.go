package service

import (
	"workflow/config"
	"workflow/src/database"

	"github.com/huylqvn/httpserver/log"
	"github.com/sirupsen/logrus"
)

type Service struct {
	Cfg *config.Config
	Log *logrus.Logger
	DB  database.Database
}

func NewService(cfg *config.Config, db database.Database) *Service {
	return &Service{
		Cfg: cfg,
		Log: log.Get(),
		DB:  db,
	}
}
