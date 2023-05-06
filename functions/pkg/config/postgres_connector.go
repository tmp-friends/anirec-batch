/*
postgresqlと接続するクラス
*/
package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const DRIVER_NAME = "postgres"

type PostgresConnector struct {
	Conn *sql.DB
}

func NewPostgresConnector() *PostgresConnector {
	conf := LoadDBConfig()
	dsn := createDSN(*conf)

	conn, err := sql.Open(DRIVER_NAME, dsn)
	if err != nil {
		panic(err)
	}

	return &PostgresConnector{
		Conn: conn,
	}
}

func createDSN(postgresInfo PostgresInfo) string {
	dataSourceName := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		postgresInfo.User,
		postgresInfo.Pass,
		postgresInfo.Host,
		postgresInfo.Port,
		postgresInfo.DBName,
	)

	return dataSourceName
}
