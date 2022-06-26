package config

import (
	"database/sql"
	"go.uber.org/zap"
)

var (
	Db     *sql.DB
	Logger *zap.Logger
)
