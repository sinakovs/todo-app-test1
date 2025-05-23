package main

import (
	"net/http"
)

func routes() http.Handler {

	api := http.NewServeMux()

	api.HandleFunc("GET /todos", getTodos)
	api.HandleFunc("POST /todos", addTodo)
	api.HandleFunc("GET /todos/one", getTodoFromOneFile)
	api.HandleFunc("GET /todos/", getTodo)
	api.HandleFunc("DELETE /todos/", deleteTodo)
	api.HandleFunc("PATCH /todos/", toggleTodoStatus)
	api.HandleFunc("GET /check", getTodosFromFile)
	api.HandleFunc("POST /save", saveDataToFile)
	api.HandleFunc("POST /parallel", getTodosAndSave)
	api.HandleFunc("GET /parallel/file", getTodosFromFileParallel)
	api.HandleFunc("GET /parallel/rwm", getTodosFromFileParallelRWM)

	api.HandleFunc("GET /cache", getTodosFromCache)
	api.HandleFunc("POST /cache", addTodoToCache)
	api.HandleFunc("GET /cache/", getTodoByIDFromCache)

	return api
}
