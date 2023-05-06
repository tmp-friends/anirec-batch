package config

import "os"

type PostgresInfo struct {
	User   string
	Pass   string
	Host   string
	Port   string
	DBName string
}

// DBの設定情報を読み込む
func LoadDBConfig() *PostgresInfo {
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPass := os.Getenv("POSTGRES_PASS")
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresDBName := os.Getenv("POSTGRES_DBNAME")

	config := &PostgresInfo{
		User:   postgresUser,
		Pass:   postgresPass,
		Host:   postgresHost,
		Port:   postgresPort,
		DBName: postgresDBName,
	}

	return config
}
