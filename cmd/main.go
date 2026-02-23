package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"test-task/db"
	"test-task/internal/domain"
	"test-task/internal/handlers"
	"test-task/internal/kafka"
	middleware1 "test-task/internal/middleware"
	"test-task/internal/repository"
	"test-task/internal/service"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf(".env kal: %v", err)
	}

	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	if err := database.AutoMigrate(&domain.User{}, &domain.Note{}); err != nil {
		log.Fatal("Migration failed:", err)
	}

	kafkaConfig := kafka.NewConfig()
	kafkaProducer := kafka.NewProducer(kafkaConfig)
	defer kafkaProducer.Close()

	noteRepo := repository.NewNotesRepository(database)
	userRepo := repository.NewUserRepository(database)

	kafkaConsumer := kafka.NewConsumer(kafkaConfig, noteRepo)
	defer kafkaConsumer.Close()

	ctx := context.Background()
	go kafkaConsumer.ConsumeMessages(ctx)

	noteService := service.NewNoteService(noteRepo, kafkaProducer)
	userService := service.NewUserService(userRepo)

	noteHandlers := handlers.NewNotesHandlers(noteService)
	userHandlers := handlers.NewUserHandler(userService)

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.POST("/register", userHandlers.Register)
	e.POST("/login", userHandlers.Login)

	api := e.Group("/api")
	api.Use(middleware1.AuthMiddleware)

	api.GET("/notes", noteHandlers.GetNotes)
	api.GET("/notes/id/:id", noteHandlers.GetNoteById)
	api.POST("/notes", noteHandlers.PostNotes)
	api.DELETE("/notes/:id", noteHandlers.DeleteNote)
	api.PATCH("notes/update/:id", noteHandlers.UpdateNote)

	e.Start("localhost:8080")

}
