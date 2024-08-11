package command

import (
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	appHttp "github.com/rusinov-artem/gophermart/app/http"
	"github.com/rusinov-artem/gophermart/cmd/gophermart/config"
)

type Server interface {
	Run()
}

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "gophermart",
	}

	defaultAddress := ":80"
	if v, ok := os.LookupEnv("RUN_ADDRESS"); ok {
		defaultAddress = v
	}

	defaultDSN := ""
	if v, ok := os.LookupEnv("DATABASE_URI"); ok {
		defaultDSN = v
	}

	accrualAddress := ""
	if v, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); ok {
		accrualAddress = v
	}

	cmd.PersistentFlags().StringP("address", "a", defaultAddress, "address to listen to")
	cmd.PersistentFlags().StringP("database", "d", defaultDSN, "address to listen to")
	cmd.PersistentFlags().StringP("accrual", "r", accrualAddress, "address to listen to")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		cfg := config.New().Load(cmd)
		srv := BuildServer(cfg)
		srv.Run()
	}

	return cmd
}

var BuildServer = func(cfg *config.Config) Server {
	logger, _ := zap.NewProduction()
	logger = logger.With(zap.Any("config", cfg))

	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	s := appHttp.NewServer(cfg.Address, mux, logger)

	return s
}
