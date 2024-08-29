package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rusinov-artem/gophermart/cmd/gophermart/command"
	"github.com/rusinov-artem/gophermart/cmd/gophermart/config"
)

func Test_ConfigurationFromEnv(t *testing.T) {
	t.Setenv("RUN_ADDRESS", "gophermart.test:9988")
	t.Setenv("DATABASE_URI", "postgres")
	t.Setenv("ACCRUAL_SYSTEM_ADDRESS", "accrual-system.test:8080")
	command.BuildServer = func(cfg *config.Config) command.Server {
		assert.Equal(t, "gophermart.test:9988", cfg.Address)
		assert.Equal(t, "postgres", cfg.DatabaseDSN)
		assert.Equal(t, "accrual-system.test:8080", cfg.AccrualSystemAddress)

		return server{}
	}

	cmd := command.RootCmd()
	err := cmd.Execute()
	require.NoError(t, err)
}

func Test_ConfigurationFromCmd(t *testing.T) {
	command.BuildServer = func(cfg *config.Config) command.Server {
		assert.Equal(t, "gophermart.dev:9988", cfg.Address)
		assert.Equal(t, "mysql", cfg.DatabaseDSN)
		assert.Equal(t, "accrual-system.dev:8080", cfg.AccrualSystemAddress)

		return server{}
	}

	cmd := command.RootCmd()
	cmd.SetArgs([]string{
		"",
		"-a", "gophermart.dev:9988",
		"-d", "mysql",
		"-r", "accrual-system.dev:8080",
	})
	err := cmd.Execute()
	require.NoError(t, err)
}

type server struct {
}

func (s server) Run() {}
