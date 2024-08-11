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
	if err := goose.SetDialect("pgx"); err != nil {
		log.Error("goose: unable to set dialect", zap.Error(err))
	}

	migrationDir := os.Getenv("MIGRATION_DIR")
	if migrationDir == "" {
		dir, err := os.Getwd()
		if err != nil {
			log.Error("goose: unable get current working directory", zap.Error(err))
		}

		migrationDir = fmt.Sprintf("%s/app/migration", dir)
	}

	log.Info("goose: migrations dir", zap.String("dir", migrationDir))
	err := goose.Up(db, migrationDir)
	if err != nil {
		log.Error("goose: unable to run migration", zap.Error(err))
	}
}
