package database

import (
	"context"
	"database/sql"
	"github.com/TonimatasDEV/ReposiGO/configuration"
	"github.com/TonimatasDEV/ReposiGO/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "modernc.org/sqlite"
	"strconv"
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

func saveMySQLandMarianDB(dbConfig configuration.Database, username string, hashedToken string, writeAccess string, readAccess string) error {
	db, err := sql.Open("mysql", dbConfig.User+":"+dbConfig.Password+"@tcp("+dbConfig.Host+":"+strconv.Itoa(dbConfig.Port)+")/"+dbConfig.Name)
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

func savePostgreSQL(dbConfig configuration.Database, username string, hashedToken string, writeAccess string, readAccess string) error {
	db, err := pgxpool.New(context.Background(), "postgres://"+dbConfig.User+":"+dbConfig.Password+"@"+dbConfig.Host+":"+strconv.Itoa(dbConfig.Port)+"/"+dbConfig.Name)
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

func saveMongoDB(dbConfig configuration.Database, username string, hashedToken string, writeAccess string, readAccess string) error {
	clientOptions := options.Client().ApplyURI("mongodb://" + dbConfig.User + ":" + dbConfig.Password + "@" + dbConfig.Host + ":" + strconv.Itoa(dbConfig.Port))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	defer utils.MongoDBDisconnectError(client)

	collection := client.Database(dbConfig.Name).Collection("sessions")

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
