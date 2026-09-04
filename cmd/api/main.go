package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
	"task-api/internal/middleware"
	"task-api/internal/handler"
	"task-api/internal/repository"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")

	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	taskRepo := repository.NewTaskRepository(dbpool)
	taskHandler := handler.NewTaskHandler(taskRepo)

	mux := http.NewServeMux()

    mux.HandleFunc("POST /tasks", taskHandler.CreateTask)
	mux.HandleFunc("GET /tasks", taskHandler.GetTasks)
	mux.HandleFunc("GET /tasks/{id}", taskHandler.GetTaskByID)
	mux.HandleFunc("PUT /tasks/{id}", taskHandler.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTask)

	handlerWithMiddleware := middleware.Logger(middleware.Recover(mux))

	log.Println("Server running on http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", handlerWithMiddleware))
}

// QUESTION 1: Why do we have to close the resp.Body

// Suppose the server sends you 1 MB of data.
// Your program might read it like this:
// body, err := io.ReadAll(resp.Body)
// While you're doing this, the HTTP client may have a network connection and other resources associated with that response.
// When you're finished:
// resp.Body.Close()
// tells the HTTP client:
// "I'm done with this response body. Release whatever resources are associated with it."
// That's why the documentation says:
// defer resp.Body.Close()
// also this resp.Body recieves a stream of data
// resp is a *http.Response struct.