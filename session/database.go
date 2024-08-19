package session

import (
	"errors"
	"github.com/TonimatasDEV/ReposiGO/configuration"
	"github.com/TonimatasDEV/ReposiGO/session/database"
	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

func deleteSession(username string, hashedToken string, writeAccess string, readAccess string) error {
	return nil
}

func saveSession(username string, hashedToken string, writeAccess string, readAccess string) error {
	dbConfig := configuration.ServerConfig.Database

	switch dbConfig.Type {
	case "sqlite":
		return database.SaveSQLite(username, hashedToken, writeAccess, readAccess)
	case "mysql", "mariadb":
		return database.SaveMySQLandMarianDB(dbConfig, username, hashedToken, writeAccess, readAccess)
	case "postgresql":
		return database.SavePostgresql(dbConfig, username, hashedToken, writeAccess, readAccess)
	case "mongodb":
		return database.SaveMongoDB(dbConfig, username, hashedToken, writeAccess, readAccess)
	default:
		return errors.New("invalid database type")
	}
}

func readSessions() (map[string]Session, error) {
	dbConfig := configuration.ServerConfig.Database

	switch dbConfig.Type {
	case "sqlite":
		return database.ReadSessionsSQLite()
	case "mysql", "mariadb":
		return database.ReadSessionsMySQLandMariaDB(dbConfig)
	case "postgresql":
		return database.ReadSessionsPostgresql(dbConfig)
	case "mongodb":
		return database.ReadSessionsMongoDB(dbConfig)
	default:
		return nil, errors.New("invalid database type")
	}
}
