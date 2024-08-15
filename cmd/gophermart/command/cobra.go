package command

import (
	"context"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/rusinov-artem/gophermart/app/action/balance/get"
	"github.com/rusinov-artem/gophermart/app/action/login"
	"github.com/rusinov-artem/gophermart/app/action/order/add"
	"github.com/rusinov-artem/gophermart/app/action/order/list"
	"github.com/rusinov-artem/gophermart/app/action/register"
	getWithdrawalsAction "github.com/rusinov-artem/gophermart/app/action/withdraw/get"
	"github.com/rusinov-artem/gophermart/app/action/withdraw/process"
	"github.com/rusinov-artem/gophermart/app/crypto"
	appHttp "github.com/rusinov-artem/gophermart/app/http"
	appHandler "github.com/rusinov-artem/gophermart/app/http/handler"
	"github.com/rusinov-artem/gophermart/app/http/middleware"
	appRouter "github.com/rusinov-artem/gophermart/app/http/router"
	"github.com/rusinov-artem/gophermart/app/migration"
	"github.com/rusinov-artem/gophermart/app/service/accrual"
	"github.com/rusinov-artem/gophermart/app/service/accrual/client"
	"github.com/rusinov-artem/gophermart/app/service/auth"
	"github.com/rusinov-artem/gophermart/app/service/order"
	withdrawService "github.com/rusinov-artem/gophermart/app/service/withdraw"
	"github.com/rusinov-artem/gophermart/app/storage"
	"github.com/rusinov-artem/gophermart/cmd/gophermart/config"
)

type Server interface {
	Run()
}

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "gophermart",
	}

	cmd.AddCommand(migrate())

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

	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		panic(err)
	}

	migration.Migrate(logger, dbpool)

	c := chi.NewRouter()
	c.Use(middleware.Logger(logger))

	handler := appHandler.New()

	handler.RegisterAction = func(ctx context.Context) appHandler.RegisterAction {
		s := storage.NewRegistrationStorage(ctx, dbpool)

		return register.New(s, logger, crypto.NewTokenGenerator())
	}

	handler.LoginAction = func(ctx context.Context) appHandler.LoginAction {
		s := storage.NewRegistrationStorage(ctx, dbpool)

		h := login.New(s, logger, crypto.NewTokenGenerator())
		h.CheckPasswordHash = crypto.CheckPasswordHash

		return h
	}

	handler.AuthService = func(ctx context.Context) appHandler.AuthService {
		s := storage.NewRegistrationStorage(ctx, dbpool)
		return auth.NewService(s)
	}

	handler.AddOrderAction = func(ctx context.Context) appHandler.AddOrderAction {
		s := storage.NewRegistrationStorage(ctx, dbpool)
		return add.New(s, logger)
	}

	handler.ListOrdersAction = func(ctx context.Context) appHandler.ListOrdersAction {
		storage := storage.NewRegistrationStorage(ctx, dbpool)

		accrualClient := client.New(ctx, cfg.AccrualSystemAddress)
		accrualService := accrual.NewService(accrualClient, storage, logger)
		orderService := order.NewOrderService(logger, storage, accrualService)

		return list.New(orderService)
	}

	handler.GetBalanceAction = func(ctx context.Context) appHandler.GetBalanceAction {
		storage := storage.NewRegistrationStorage(ctx, dbpool)

		accrualClient := client.New(ctx, cfg.AccrualSystemAddress)
		accrualService := accrual.NewService(accrualClient, storage, logger)
		orderService := order.NewOrderService(logger, storage, accrualService)

		return get.New(orderService)
	}

	handler.WithdrawAction = func(ctx context.Context) appHandler.WithdrawAction {
		s := storage.NewRegistrationStorage(ctx, dbpool)
		accrualClient := client.New(ctx, cfg.AccrualSystemAddress)
		accrualService := accrual.NewService(accrualClient, s, logger)
		orderService := order.NewOrderService(logger, s, accrualService)

		txFactory := func(login string) withdrawService.Transaction {
			return storage.NewWithdrawTx(ctx, dbpool, login)
		}

		withdrawService := withdrawService.NewWithdrawService(txFactory, logger)

		return process.New(orderService, withdrawService)
	}

	handler.GetWithdrawalsAction = func(ctx context.Context) appHandler.GetWithdrawalsAction {
		s := storage.NewRegistrationStorage(ctx, dbpool)

		return getWithdrawalsAction.New(logger, s)
	}

	router := appRouter.New(c).SetHandler(handler)

	s := appHttp.NewServer(cfg.Address, router.Mux(), logger)

	return s
}
