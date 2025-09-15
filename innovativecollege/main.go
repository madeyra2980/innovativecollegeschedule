package main

import (
	"log"
	"os"

	"innovativecollege/internal/config"
	"innovativecollege/internal/database"
	"innovativecollege/internal/handlers"
	"innovativecollege/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load("config.env"); err != nil {
		log.Println("Файл config.env не найден, используем переменные окружения")
	}

	// Инициализируем конфигурацию
	cfg := config.Load()

	// Подключаемся к MongoDB
	db, err := database.Connect(cfg.MongoURI, cfg.DatabaseName)
	if err != nil {
		log.Fatal("Ошибка подключения к MongoDB:", err)
	}
	defer database.Close(database.Client)

	// Создаем коллекции
	database.CreateCollections(db)

	// Инициализируем обработчики
	h := handlers.New(db)

	// Настраиваем роуты
	r := gin.Default()
	routes.SetupRoutes(r, h)

	// Запускаем сервер
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(r.Run(":" + port))
}
