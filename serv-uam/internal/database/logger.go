package database

import (
	"log"
	"os"
	"time"

	"gorm.io/gorm/logger"
)

// gormLogger only logs query errors to the console — "record not found" is
// excluded since it's an expected outcome, not a failure. Successful
// queries stay silent (no per-query logging).
var gormLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags),
	logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  logger.Error,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
	},
)
