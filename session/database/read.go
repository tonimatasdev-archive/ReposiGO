package database

import (
	"context"
	"database/sql"
	"github.com/TonimatasDEV/ReposiGO/configuration"
	"github.com/TonimatasDEV/ReposiGO/session"
	"github.com/TonimatasDEV/ReposiGO/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "modernc.org/sqlite"
	"strconv"
	"strings"
)

//goland:noinspection SqlNoDataSourceInspection
func ReadSessionsSQLite() (map[string]session.Session, error) {
	db, err := sql.Open("sqlite", "file:sessions.db")
	if err != nil {
		return nil, err
	}

	defer utils.CloseDBError(db)

	rows, err := db.Query(`SELECT username, token_hash, write_access, read_access FROM sessions`)
	if err != nil {
		return nil, err
	}

	defer utils.CloseRowError(rows)

	var sessions = make(map[string]session.Session)

	for rows.Next() {
		var username, tokenHash, writeAccessStr, readAccessStr string

		if err := rows.Scan(&username, &tokenHash, &writeAccessStr, &readAccessStr); err != nil {
			return nil, err
		}

		rawSession := session.Session{
			Username:    username,
			HashedToken: tokenHash,
			WriteAccess: strings.Split(writeAccessStr, ","),
			ReadAccess:  strings.Split(readAccessStr, ","),
		}

		sessions[rawSession.Username] = rawSession
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

//goland:noinspection SqlNoDataSourceInspection
func ReadSessionsMySQLandMariaDB(dbConfig configuration.Database) (map[string]session.Session, error) {
	db, err := sql.Open("mysql", dbConfig.User+":"+dbConfig.Password+"@tcp("+dbConfig.Host+":"+strconv.Itoa(dbConfig.Port)+")/"+dbConfig.Name)
	if err != nil {
		return nil, err
	}

	defer utils.CloseDBError(db)

	rows, err := db.Query(`SELECT username, token_hash, write_access, read_access FROM sessions`)
	if err != nil {
		return nil, err
	}

	defer utils.CloseRowError(rows)

	var sessions = make(map[string]session.Session)

	for rows.Next() {
		var username, tokenHash, writeAccessStr, readAccessStr string

		if err := rows.Scan(&username, &tokenHash, &writeAccessStr, &readAccessStr); err != nil {
			return nil, err
		}

		rawSession := session.Session{
			Username:    username,
			HashedToken: tokenHash,
			WriteAccess: strings.Split(writeAccessStr, ","),
			ReadAccess:  strings.Split(readAccessStr, ","),
		}

		sessions[rawSession.Username] = rawSession
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

//goland:noinspection SqlNoDataSourceInspection
func ReadSessionsPostgresql(dbConfig configuration.Database) (map[string]session.Session, error) {
	db, err := pgxpool.New(context.Background(), "postgres://"+dbConfig.User+":"+dbConfig.Password+"@"+dbConfig.Host+":"+strconv.Itoa(dbConfig.Port)+"/"+dbConfig.Name)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(context.Background(), `SELECT username, token_hash, write_access, read_access FROM sessions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions = make(map[string]session.Session)

	for rows.Next() {
		var rawSession session.Session
		if err := rows.Scan(&rawSession.Username, &rawSession.HashedToken, &rawSession.WriteAccess, &rawSession.ReadAccess); err != nil {
			return nil, err
		}
		sessions[rawSession.Username] = rawSession
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

func ReadSessionsMongoDB(dbConfig configuration.Database) (map[string]session.Session, error) {
	clientOptions := options.Client().ApplyURI("mongodb://" + dbConfig.User + ":" + dbConfig.Password + "@" + dbConfig.Host + ":" + strconv.Itoa(dbConfig.Port))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	defer utils.MongoDBDisconnectError(client)

	collection := client.Database(dbConfig.Name).Collection("sessions")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	defer utils.CloseCursorError(cursor, context.TODO())

	var sessions = make(map[string]session.Session)

	for cursor.Next(context.TODO()) {
		var mongoSession MongoSession
		if err := cursor.Decode(&mongoSession); err != nil {
			return nil, err
		}

		rawSession := session.Session{
			Username:    mongoSession.Username,
			HashedToken: mongoSession.TokenHash,
			WriteAccess: strings.Split(mongoSession.WriteAccess, ","),
			ReadAccess:  strings.Split(mongoSession.ReadAccess, ","),
		}

		sessions[rawSession.Username] = rawSession
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}
