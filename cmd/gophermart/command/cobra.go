package command

import (
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	appHttp "github.com/rusinov-artem/gophermart/app/http"
	appHandler "github.com/rusinov-artem/gophermart/app/http/handler"
	appRouter "github.com/rusinov-artem/gophermart/app/http/router"
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

	c := chi.NewRouter()
	handler := appHandler.New()
	router := appRouter.New(c).SetHandler(handler)

	s := appHttp.NewServer(cfg.Address, router.Mux(), logger)

	return s
}
