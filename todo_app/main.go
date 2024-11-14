package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/DaniilStelmakh/go_final_project_main/internal/http-server/handlers"
	tasks "github.com/DaniilStelmakh/go_final_project_main/internal/service"
	sqlite "github.com/DaniilStelmakh/go_final_project_main/storage/database"
	"github.com/go-chi/chi/v5"
)

const (
	defaultWebDir = "./web/"
	defaultDBFile = "./scheduler.db"
	defaultPort   = 7540
)

func main() {

	dbFile := os.Getenv("TODO_DBFILE")

	if dbFile == "" {
		dbFile = defaultDBFile
	}
	store, err := sqlite.CreateTable(dbFile)
	if err != nil {
		log.Fatalf("error starting db: %v", err)
		return
	}
	defer func() {
		err = store.Close()
		if err != nil {
			log.Fatalf("DB failed to shutdown: %v\n", err)
		}
	}()

	portStr := os.Getenv("TODO_PORT")
	var port int

	if portStr == "" {
		port = defaultPort
	} else {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			log.Fatalf("Invalid port number: %v", err)
			return
		}
	}

	webDir := os.Getenv("TODO_WEB_DIR")

	if webDir == "" {
		webDir = defaultWebDir
	}

	taskService := tasks.New(store)

	r := chi.NewRouter()

	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	r.MethodFunc(http.MethodGet, "/api/nextdate", handlers.NextTask())

	r.MethodFunc(http.MethodPost, "/api/task", handlers.AddTask(taskService))
	r.MethodFunc(http.MethodGet, "/api/task", handlers.GetTask(taskService))
	r.MethodFunc(http.MethodPut, "/api/task", handlers.UpdateTask(taskService))
	r.MethodFunc(http.MethodDelete, "/api/task", handlers.DeleteTask(taskService))

	r.MethodFunc(http.MethodGet, "/api/tasks", handlers.GetTasks(taskService))

	r.MethodFunc(http.MethodPost, "/api/task/done", handlers.DoneTask(taskService))

	// r.MethodFunc(http.MethodPost, "/api/sigin", authenticator.Authenticator)

	address := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatalf("Server failed with err: %v\n", err)
		return
	}
}
