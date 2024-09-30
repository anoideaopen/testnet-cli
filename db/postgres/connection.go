package postgres

import (
	"context"
	"fmt"

	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
)

func NewPostgresDB(
	host string,
	port string,
	user string,
	password string,
	dbName string,
) (*pg.DB, error) {
	logger.Info(
		"postgres",
		zap.String("host", host),
		zap.String("port", port),
		zap.String("user", user),
		zap.String("password", password),
		zap.String("dbName", dbName),
	)

	address := fmt.Sprintf("%s:%s", host, port)
	db := pg.Connect(&pg.Options{
		Addr:     address,
		User:     user,
		Password: password,
		Database: dbName,
	})

	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		logger.GetLogger().Error("Error ping db!")
		return nil, err
	}
	logger.GetLogger().Info("Successfully connected!")

	if err := CreateSchema(db); err != nil {
		logger.GetLogger().Error("Error create schema!")
		return nil, err
	}
	logger.GetLogger().Info("Schema created!")
	return db, nil
}
