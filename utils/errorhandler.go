package utils

import (
	"context"
	"database/sql"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
)

func FileError(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Println("Error closing file", file.Name(), ":", err)
	}
}

func CloseDBError(db *sql.DB) {
	err := db.Close()

	if err != nil {
		log.Println("Error closing the session database", err)
	}
}

func MongoDBDisconnectError(client *mongo.Client) {
	err := client.Disconnect(context.TODO())

	if err != nil {
		log.Println("Error disconnecting from MongoDB", err)
	}
}
