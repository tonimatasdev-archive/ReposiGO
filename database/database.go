package database

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

func DeleteSession(username string, token []byte, writeAccess string, readAccess string) error {
	return nil
}

func SaveSession(username string, hashedToken string, writeAccess string, readAccess string) error {
	switch "" {
	case "sqlite":
		return saveSQLite(username, hashedToken, writeAccess, readAccess)
	case "mysql", "mariadb":
		return saveMySQLandMarianDB(username, hashedToken, writeAccess, readAccess)
	case "postgresql":
		return savePostgreSQL(username, hashedToken, writeAccess, readAccess)
	case "mongodb":
		return saveMongoDB(username, hashedToken, writeAccess, readAccess)
	default:
		return errors.New("invalid database type")
	}
}

func ReadSessions() error {
	return nil
}
