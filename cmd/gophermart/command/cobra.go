package command

import (
	"net/http"
	"os"

	appHttp "github.com/rusinov-artem/gophermart/app/http"
	"github.com/rusinov-artem/gophermart/cmd/gophermart/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "gophermart",
	}

	defaultAddress := ":80"
	if v, ok := os.LookupEnv("RUN_ADDRESS"); ok {
		defaultAddress = v
	}
	cmd.PersistentFlags().StringP("address", "a", defaultAddress, "address to listen to")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		cfg := config.New().Load(cmd)
		srv := buildServer(cfg)
		srv.Run()
	}

	return cmd
}

var buildServer = func(cfg *config.Config) *appHttp.Server {
	logger, _ := zap.NewProduction()
	logger = logger.With(zap.Any("config", cfg))

	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	s := appHttp.NewServer(cfg.Address, mux, logger)

	return s
}
