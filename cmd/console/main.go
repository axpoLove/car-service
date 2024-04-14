package main

import (
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"

	"car-service/internal/config"
)

var (
	rootCmd = &cobra.Command{
		Use:   "migrate",
		Short: "run migration",
		Run: func(cmd *cobra.Command, args []string) {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			logger.Info("please specify a migration command.")
		},
	}
	upCmd = &cobra.Command{
		Use:   "up",
		Short: "run all available migrations",
		Run: func(cmd *cobra.Command, args []string) {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			cfg, err := config.NewConfig()
			if err != nil {
				logger.Error("failed to load config", "error", err)
				return
			}
			connString := cfg.Postgres.ConnString()
			db, err := goose.OpenDBWithDriver("postgres", connString)
			if err != nil {
				logger.Error("failed to connect to database", "error", err)
				return
			}
			err = goose.Up(db, cfg.Postgres.MigrationsDir)
			if err != nil {
				logger.Error("failed to up migrations", "error", err)
				return
			}
		},
	}
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "create a new migration file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			cfg, err := config.NewConfig()
			if err != nil {
				logger.Error("failed to init config", "error", err)
				return
			}
			connString := cfg.Postgres.ConnString()
			db, err := goose.OpenDBWithDriver("postgres", connString)
			if err != nil {
				logger.Error("failed to connect to database", "error", err)
				return
			}
			name := args[0]
			err = goose.Create(db, cfg.Postgres.MigrationsDir, name, "sql")
			if err != nil {
				logger.Error("failed to create the migration", "error", err)
				return
			}
			logger.Info("a new migration created")
		},
	}
)

func init() {
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(createCmd)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := rootCmd.Execute(); err != nil {
		logger.Error("failed to execute root cmd", "error", err)
		os.Exit(1)
	}
}
