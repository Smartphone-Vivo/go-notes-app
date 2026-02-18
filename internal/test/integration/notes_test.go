package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"test-task/internal/domain"
	"test-task/internal/handlers"
	"test-task/internal/middleware"
	"test-task/internal/repository"
	"test-task/internal/service"
)

var (
	testDB    *gorm.DB
	cleanup   func()
	testApp   *echo.Echo
	testToken string
)

func TestMain(m *testing.M) {
	setupTestContainer()
	setupTestUser()

	code := m.Run()

	if cleanup != nil {
		cleanup()
	}

	os.Exit(code)
}

func setupTestContainer() {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		panic(err)
	}

	host, _ := pgContainer.Host(ctx)
	port, _ := pgContainer.MappedPort(ctx, "5432")

	connStr := "host=" + host + " port=" + port.Port() + " user=test password=test dbname=testdb sslmode=disable"

	testDB, err = gorm.Open(gormpostgres.Open(connStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	testDB.AutoMigrate(&domain.User{}, &domain.Note{})

	noteRepo := repository.NewNotesRepository(testDB)
	userRepo := repository.NewUserRepository(testDB)

	noteService := service.NewNoteService(noteRepo)
	userService := service.NewUserService(userRepo)

	noteHandlers := handlers.NewNotesHandlers(noteService)
	userHandlers := handlers.NewUserHandler(userService)

	testApp = echo.New()

	testApp.POST("/register", userHandlers.Register)
	testApp.POST("/login", userHandlers.Login)

	api := testApp.Group("/api")
	api.Use(middleware.AuthMiddleware)

	api.GET("/notes", noteHandlers.GetNotes)
	api.POST("/notes", noteHandlers.PostNotes)
	api.GET("/notes/id/:id", noteHandlers.GetNoteById)
	api.DELETE("/notes/:id", noteHandlers.DeleteNote)

	cleanup = func() { pgContainer.Terminate(ctx) }
}

func setupTestUser() {
	userJSON := `{"username":"testuser","password":"testpass"}`
	req := httptest.NewRequest("POST", "/register", bytes.NewReader([]byte(userJSON)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testApp.ServeHTTP(rec, req)

	var resp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)

	testToken = resp["token"].(string)
}

func cleanTable() {
	testDB.Exec("DELETE FROM notes")
}

func TestCreateAndGetNote(t *testing.T) {
	cleanTable()

	noteJSON := `{"title":"Тест","content":"Привет"}`
	req := httptest.NewRequest("POST", "/api/notes", bytes.NewReader([]byte(noteJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken)
	rec := httptest.NewRecorder()
	testApp.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)

	var created map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &created)
	id := created["id"].(string)

	req = httptest.NewRequest("GET", "/api/notes/id/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	rec = httptest.NewRecorder()
	testApp.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestDeleteNote(t *testing.T) {
	cleanTable()

	noteJSON := `{"title":"Тест","content":"Привет"}`
	req := httptest.NewRequest("POST", "/api/notes", bytes.NewReader([]byte(noteJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken)
	rec := httptest.NewRecorder()
	testApp.ServeHTTP(rec, req)

	var created map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &created)
	id := created["id"].(string)
	assert.NotEmpty(t, id)

	req = httptest.NewRequest("DELETE", "/api/notes/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	rec = httptest.NewRecorder()
	testApp.ServeHTTP(rec, req)
	assert.Equal(t, 204, rec.Code)

	req = httptest.NewRequest("GET", "/api/notes/id/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	rec = httptest.NewRecorder()
	testApp.ServeHTTP(rec, req)
	assert.Equal(t, 404, rec.Code)
}
