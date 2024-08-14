package migration

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

func Migrate(log *zap.Logger, pool *pgxpool.Pool) {
	db := stdlib.OpenDBFromPool(pool)
	_ = goose.SetDialect("pgx")

	migrationDir := os.Getenv("MIGRATION_DIR")
	if migrationDir == "" {
		dir, _ := os.Getwd()

		migrationDir = fmt.Sprintf("%s/app/migration", dir)
	}

	log.Info("goose: migrations dir", zap.String("dir", migrationDir))
	err := goose.Up(db, migrationDir)
	if err != nil {
		log.Error("goose: unable to run migration", zap.Error(err))
	}
}
