package database

import (
	"errors"
	"github.com/TonimatasDEV/ReposiGO/configuration"
	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

func DeleteSession(username string, hashedToken string, writeAccess string, readAccess string) error {
	return nil
}

func SaveSession(username string, hashedToken string, writeAccess string, readAccess string) error {
	dbConfig := configuration.ServerConfig.Database

	switch dbConfig.Type {
	case "sqlite":
		return saveSQLite(username, hashedToken, writeAccess, readAccess)
	case "mysql", "mariadb":
		return saveMySQLandMarianDB(dbConfig, username, hashedToken, writeAccess, readAccess)
	case "postgresql":
		return savePostgreSQL(dbConfig, username, hashedToken, writeAccess, readAccess)
	case "mongodb":
		return saveMongoDB(dbConfig, username, hashedToken, writeAccess, readAccess)
	default:
		return errors.New("invalid database type")
	}
}

func ReadSessions() error {
	return nil
}
