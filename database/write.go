package database

import (
	"context"
	"database/sql"
	"github.com/TonimatasDEV/ReposiGO/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "modernc.org/sqlite"
)

func saveSQLite(username string, hashedToken string, writeAccess string, readAccess string) error {
	db, err := sql.Open("sqlite", "file:sessions.db")
	if err != nil {
		return err
	}

	defer utils.CloseDBError(db)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		token_hash TEXT NOT NULL,
		write_access TEXT NOT NULL,
		read_access TEXT NOT NULL
	);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO sessions (username, token_hash, write_access, read_access) VALUES (?, ?)`, username, hashedToken, writeAccess, readAccess)
	if err != nil {
		return err
	}

	return nil
}

func saveMySQLandMarianDB(username string, hashedToken string, writeAccess string, readAccess string) error {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/dbname") // TODO: Add config
	if err != nil {
		return err
	}

	defer utils.CloseDBError(db)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTO_INCREMENT,
		username TEXT NOT NULL,
		token_hash TEXT NOT NULL,
		write_access TEXT NOT NULL,
		read_access TEXT NOT NULL
	);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO sessions (username, token_hash, write_access, read_access) VALUES (?, ?)`, username, hashedToken, writeAccess, readAccess)
	if err != nil {
		return err
	}

	return nil
}

func savePostgreSQL(username string, hashedToken string, writeAccess string, readAccess string) error {
	db, err := pgxpool.New(context.Background(), "postgres://username:password@localhost:5432/mydb") // TODO: Add config
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS sessions (
		id SERIAL PRIMARY KEY,
		username TEXT NOT NULL,
		token_hash TEXT NOT NULL,
		write_access TEXT NOT NULL,
		read_access TEXT NOT NULL
	);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), `INSERT INTO sessions (username, token_hash, write_access, read_access) VALUES ($1, $2)`, username, hashedToken, writeAccess, readAccess)
	if err != nil {
		return err
	}

	return nil
}

type MongoSession struct {
	Username    string `bson:"username"`
	TokenHash   string `bson:"token_hash"`
	WriteAccess string `bson:"write_access"`
	ReadAccess  string `bson:"read_access"`
}

func saveMongoDB(username string, hashedToken string, writeAccess string, readAccess string) error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // TODO: Add config
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	defer utils.MongoDBDisconnectError(client)

	collection := client.Database("mydb").Collection("sessions") // TODO: Add config

	mongoSession := MongoSession{
		Username:    username,
		TokenHash:   hashedToken,
		WriteAccess: writeAccess,
		ReadAccess:  readAccess,
	}

	_, err = collection.InsertOne(context.TODO(), mongoSession)
	if err != nil {
		return err
	}

	return nil
}
