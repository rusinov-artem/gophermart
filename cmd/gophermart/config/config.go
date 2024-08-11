package config

import (
	"github.com/spf13/cobra"
)

type Config struct {
	Address              string
	DatabaseDSN          string
	AccrualSystemAddress string
}

func New() *Config {
	return &Config{}
}

func (c *Config) Load(cmd *cobra.Command) *Config {
	c.Address, _ = cmd.Flags().GetString("address")
	c.DatabaseDSN, _ = cmd.Flags().GetString("database")
	c.AccrualSystemAddress, _ = cmd.Flags().GetString("accrual")
	return c
}
