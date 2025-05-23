package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type todo struct {
	TodosNumber int    `json;"id"`
	ID          string `json;"id"`
	Item        string `json;"item"`
	Completed   bool   `json;"completed"`
}

type Task struct {
	ID   int     `json:"id"`
	Todo *[]todo `json:"title"`
	Done bool    `json:"done"`
}

var taskQueue = make(chan Task, 100)
var mu sync.Mutex

func getTodoDataFromOneFile() (*[]todo, error) {
	var todosForTest = []todo{}
	file, err := os.Open(dataPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := todoDecode(scanner.Text())
		//fmt.Println(line)
		todosForTest = append(todosForTest, line)
	}

	return &todosForTest, nil
}

func getTodoDataFromFile() (*[]todo, error) {
	var testTodos = []todo{}

	file, err := os.Open(dataPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := todoDecode(scanner.Text())
		//fmt.Println(line)
		testTodos = append(testTodos, line)
	}

	file2, err := os.Open(data2Path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file2.Close()

	scanner = bufio.NewScanner(file2)
	for scanner.Scan() {
		line := todoDecode(scanner.Text())
		//fmt.Println(line)
		testTodos = append(testTodos, line)
	}

	file3, err := os.Open(data3Path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file3.Close()

	scanner = bufio.NewScanner(file3)
	for scanner.Scan() {
		line := todoDecode(scanner.Text())
		//fmt.Println(line)
		testTodos = append(testTodos, line)
	}

	return &testTodos, nil
}

func writeDataInFile(filepath string, todos *[]todo) error {

	_, err := os.Stat(filepath)
	if err != nil {
		fmt.Println(err.Error())
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	byteSlice := todoEncode(todos)

	_, err = file.Write(byteSlice)
	if err != nil {
		return err
	}

	return nil
}

func writeDataInFileParallel(task Task) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := os.Stat(dataPath)
	if err != nil {
		fmt.Println(err.Error())
	}

	file, err := os.Create(dataPath)
	if err != nil {
		return err
	}
	defer file.Close()

	byteSlice := todoEncode(task.Todo)

	_, err = file.Write(byteSlice)
	if err != nil {
		return err
	}

	//task.Done = true
	//fmt.Printf("task %d is done %v", task.ID, task.Done)

	return nil
}

func todoEncode(todos *[]todo) []byte {

	str := ""

	for _, t := range *todos {
		str += fmt.Sprintf("%d::%s::%s::%t\n", t.TodosNumber, t.ID, t.Item, t.Completed)
	}

	byteSlice := []byte(str)

	return byteSlice
}

func todoDecode(line string) todo {
	var testTodo todo

	myString := strings.Split(string(line), "::")
	completed, err := strconv.ParseBool(myString[3])
	if err != nil {
		fmt.Println("Error parsing str to bool:", err)
		return testTodo
	}
	number, err := strconv.Atoi(myString[0])
	if err != nil {
		fmt.Println("Error parsing str to int:", err)
		return testTodo
	}
	testTodo.TodosNumber = number
	testTodo.ID = myString[1]
	testTodo.Item = myString[2]
	testTodo.Completed = completed

	return testTodo
}

func getTodoByID(id string) (*todo, error) {
	for i, t := range *todos {
		if t.ID == id {
			return &(*todos)[i], nil
		}
	}

	return nil, errors.New("todo not found")
}

func getIndexById(id string) (int, error) {
	for i, t := range *todos {
		if t.ID == id {
			return i, nil
		}
	}

	return -1, errors.New("todo not found")
}

func getTodosLastNumber(todos *[]todo) int {
	number := 0
	if len(*todos) > 0 {
		number = (*todos)[len(*todos)-1].TodosNumber
	}
	return number
}

func startWorkerPool(n int) {
	for i := 0; i < n; i++ {
		go func(workerID int) {
			for task := range taskQueue {
				fmt.Printf("[Worker %d] Processing the task: %v\n", workerID, task)
				writeDataInFileParallel(task)
			}
		}(i)
	}
}

func combine(data1 *[]todo, data2 *[]todo, data3 *[]todo) *[]todo {
	result := append(append((*data1), (*data2)...), (*data3)...)

	return &result
}

func getTodoDataFromFileP(n int) (*[]todo, error) {
	var testTodos = []todo{}
	var npath string
	switch n {
	case 1:
		npath = dataPath
	case 2:
		npath = data2Path
	case 3:
		npath = data3Path
	default:
		fmt.Println("There is no such option in func getTodoDataFromFileP()")
		return &testTodos, errors.New("wrong function getTodoDataFromFileP() usage")
	}

	file, err := os.Open(npath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := todoDecode(scanner.Text())
		//fmt.Println(line)
		testTodos = append(testTodos, line)
	}
	return &testTodos, nil
}
