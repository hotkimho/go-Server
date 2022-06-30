package config

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	Db     *sql.DB
	Logger *zap.Logger
	Rdb    *redis.Client
)
