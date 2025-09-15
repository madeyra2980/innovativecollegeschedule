package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect(uri, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	Client = client

	// Проверяем подключение
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Успешно подключились к MongoDB")
	return client.Database(dbName), nil
}

func Close(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Println("Ошибка при отключении от MongoDB:", err)
	}
}

func CreateCollections(db *mongo.Database) {
	collections := []string{"groups", "students", "teachers", "schedules", "subjects", "lessons", "time_slots"}

	for _, collectionName := range collections {
		// Создаем коллекцию если она не существует
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := db.CreateCollection(ctx, collectionName)
		if err != nil {
			// Коллекция уже существует, это нормально
			log.Printf("Коллекция %s уже существует или создана", collectionName)
		} else {
			log.Printf("Создана коллекция %s", collectionName)
		}
	}
}
