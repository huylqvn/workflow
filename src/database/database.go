package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
	"workflow/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database interface {
	DB() (*sql.DB, error)
	GetDB() *gorm.DB
	Transaction() (tx *gorm.DB, commit func(err error))
}

func New(cfg *config.Config) (Database, error) {
	con, err := NewDB(cfg)
	if err != nil {
		return nil, err
	}
	return &db{con: con, driver: cfg.Database.Type}, err
}

func GormLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             10 * time.Second, // Slow SQL threshold
			LogLevel:                  logger.Silent,    // Log level
			IgnoreRecordNotFoundError: true,             // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,            // Disable color
		},
	)
}

func NewDB(cfg *config.Config) (*gorm.DB, error) {
	return NewPG(cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName, cfg.Database.SSL)
}

func NewPG(user, pass, host, port, dbName, ssl string) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, pass,
		host, port,
		dbName,
		ssl,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: GormLogger(),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(4000)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, nil
}

type db struct {
	driver string
	con    *gorm.DB
}

func (c *db) DB() (*sql.DB, error) {
	return c.con.DB()
}

func (c *db) GetDB() *gorm.DB {
	return c.con
}

func (c *db) Transaction() (tx *gorm.DB, commit func(err error)) {
	tx = c.con.Begin()

	commit = func(err error) {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}

	return tx, commit
}
