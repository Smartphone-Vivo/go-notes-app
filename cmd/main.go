package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"test-task/db"
	"test-task/internal/domain"
	"test-task/internal/handlers"
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

	if err := database.AutoMigrate(&domain.Note{}); err != nil {
		log.Fatal("Migration failed:", err)
	}

	e := echo.New()

	noteRepo := repository.NewNotesRepository(database)
	noteService := service.NewNoteService(noteRepo)
	noteHandlers := handlers.NewNotesHandlers(noteService)

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/notes", noteHandlers.GetNotes)
	e.GET("/notes/id/:id", noteHandlers.GetNoteById)
	e.POST("/notes", noteHandlers.PostNotes)
	e.DELETE("/notes/:id", noteHandlers.DeleteNote)
	e.PATCH("notes/update/:id", noteHandlers.UpdateNote)

	e.Start("localhost:8080")

}
