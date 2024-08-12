package command

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/rusinov-artem/gophermart/app/migration"
)

func migrate() *cobra.Command {
	cmd := &cobra.Command{
		Use: "migrate",
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewProduction()
		dsn, _ := cmd.Flags().GetString("database")
		pool, err := pgxpool.New(context.Background(), dsn)
		if err != nil {
			logger.Fatal("unable to connect to database", zap.Error(err))
		}

		migration.Migrate(logger, pool)
	}

	return cmd
}
