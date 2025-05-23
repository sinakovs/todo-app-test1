package main

import (
	"net/http"
	"time"
)

var todos, _ = getTodoDataFromFile()
var currentTodosInFile = len(*todos)

var todoCache = NewTodoCache()

func main() {

	data, _ := getTodoDataFromFile()

	for _, t := range *data {
		todoCache.Set(t)
	}

	startWorkerPool(10)

	port := "8080"

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        routes(),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	server.ListenAndServe()
}
