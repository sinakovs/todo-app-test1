package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var dataPath = "./data/data.txt"
var data2Path = "./data/data2.txt"
var data3Path = "./data/data3.txt"

func getTodos(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}

}

func getTodoFromOneFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tempTodo, err := getTodoDataFromOneFile()

	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(tempTodo); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo todo

	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	//newTodo.TodosNumber = getTodosLastNumber(todos) + 1
	newTodo.TodosNumber = len(*todos) + 1

	//fmt.Printf("Received: \n TodosNumber = %d \n Id = %s \n Item = %s \n Completed = %t\n", newTodo.TodosNumber, newTodo.ID, newTodo.Item, newTodo.Completed)
	*todos = append(*todos, newTodo)

	if newTodo.TodosNumber-currentTodosInFile > 5 {
		err := writeDataInFile(dataPath, todos)
		if err != nil {
			http.Error(w, "Troubles with writing data in file", http.StatusInternalServerError)
			fmt.Println(err.Error())
			return
		}
		currentTodosInFile = newTodo.TodosNumber
	}

}

func getTodo(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/todos/")
	if path == "" || strings.Contains(path, "/") {
		http.Error(w, "INvalid ID format", http.StatusBadRequest)
		return
	}

	todo, err := getTodoByID(path)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Received ID: %s\n", path)

}

func toggleTodoStatus(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/todos/")
	if path == "" || strings.Contains(path, "/") {
		http.Error(w, "INvalid ID format", http.StatusBadRequest)
		return
	}

	todo, err := getTodoByID(path)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo.Completed = !todo.Completed

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	//w.WriteHeader(http.StatusOK)
}

func getTodosFromFile(w http.ResponseWriter, r *http.Request) {

	todosTest, err := getTodoDataFromFile()
	if err != nil {
		http.Error(w, "Error getting data form file response", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(todosTest); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}

}

func deleteTodo(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/todos/")
	if path == "" || strings.Contains(path, "/") {
		http.Error(w, "INvalid ID format", http.StatusBadRequest)
		return
	}

	index, err := getIndexById(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	*todos = append((*todos)[:index], (*todos)[index+1:]...)

}

func saveDataToFile(w http.ResponseWriter, r *http.Request) {
	err := writeDataInFile(dataPath, todos)
	if err != nil {
		http.Error(w, "Troubles with writing data in file", http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getTodosAndSave(w http.ResponseWriter, r *http.Request) {
	var newTodo todo
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	newTodo.TodosNumber = getTodosLastNumber(todos) + 1
	//fmt.Printf("Received: \n TodosNumber = %d \n Id = %s \n Item = %s \n Completed = %t\n", newTodo.TodosNumber, newTodo.ID, newTodo.Item, newTodo.Completed)
	*todos = append(*todos, newTodo)
	task.ID = newTodo.TodosNumber
	task.Todo = todos
	task.Done = false

	if newTodo.TodosNumber-currentTodosInFile > 5 {
		taskQueue <- task
		currentTodosInFile = newTodo.TodosNumber
	}
}

func getTodosFromFileParallel(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	ch1 := make(chan *[]todo)
	ch2 := make(chan *[]todo)
	ch3 := make(chan *[]todo)

	wg.Add(3)

	go func() {
		defer wg.Done()
		todoTest, err := getTodoDataFromFileP(1)
		if err != nil {
			http.Error(w, "Error getting data form file response", http.StatusInternalServerError)
			ch1 <- nil
			return
		}
		ch1 <- todoTest
	}()

	go func() {
		defer wg.Done()
		todoTest, err := getTodoDataFromFileP(2)
		if err != nil {
			http.Error(w, "Error getting data form file response", http.StatusInternalServerError)
			ch2 <- nil
			return
		}
		ch2 <- todoTest
	}()

	go func() {
		defer wg.Done()
		todoTest, err := getTodoDataFromFileP(3)
		if err != nil {
			http.Error(w, "Error getting data form file response", http.StatusInternalServerError)
			ch3 <- nil
			return
		}
		ch3 <- todoTest
	}()

	data1 := <-ch1
	data2 := <-ch2
	data3 := <-ch3

	wg.Wait()
	close(ch1)
	close(ch2)
	close(ch3)

	result := combine(data1, data2, data3)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}

}

func getTodosFromFileParallelRWM(w http.ResponseWriter, r *http.Request) {
	var result []todo
	var mu sync.RWMutex
	var wg sync.WaitGroup

	collect := func(n int) {
		defer wg.Done()
		mu.RLock()
		todoList, err := getTodoDataFromFileP(n)
		mu.RUnlock()
		if err != nil {
			fmt.Println("Error getting data from file", n, ":", err)
			return
		}
		mu.Lock()
		result = append(result, *todoList...)
		mu.Unlock()
	}

	wg.Add(3)
	go collect(1)
	go collect(2)
	go collect(3)

	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func getTodosFromCache(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todos := todoCache.All()
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func addTodoToCache(w http.ResponseWriter, r *http.Request) {
	var newTodo todo
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	newTodo.TodosNumber = len(todoCache.All()) + 1
	todoCache.Set(newTodo)
	w.WriteHeader(http.StatusCreated)
}

func getTodoByIDFromCache(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/todos/")
	t, ok := todoCache.Get(id)
	if !ok {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}
