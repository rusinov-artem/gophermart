package test

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

//nolint:gosec // tests credentials
const baseDSN = "postgres://gophermart:gophermart@gophermart-db:5432/%s"

func CreateTestDB(dbName string) string {
	defaultConnectionString := fmt.Sprintf(baseDSN, "postgres")
	owner := "gophermart"
	template := "gophermart"

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, defaultConnectionString)
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(ctx, dropConnections(dbName))
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(ctx, dropDatabase(dbName))
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(ctx, dropConnections(template))
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(ctx, createDatabaseFromTemplate(template, dbName, owner))
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(baseDSN, dbName)
}

func dropConnections(dbName string) string {
	return fmt.Sprintf(`
		SELECT pg_terminate_backend(pg_stat_activity.pid) 
		FROM pg_stat_activity 
		WHERE pg_stat_activity.datname = '%s' AND pid <> pg_backend_pid();`,
		dbName)
}

func dropDatabase(dbName string) string {
	return fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)
}

func createDatabaseFromTemplate(template, dbName, owner string) string {
	return fmt.Sprintf(`CREATE DATABASE %s WITH TEMPLATE %q OWNER %q`, dbName, template, owner)
}
