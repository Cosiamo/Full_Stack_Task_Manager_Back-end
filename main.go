package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Tasks struct {
	ID string `json:"id"`
	TaskName string `json:"task_name"`
	TaskDetail string `json:"task_detail"`
	Date string `json:"date"`
}

var tasks []Tasks

func allTasks() {
	task := Tasks{
		ID: "1",
		TaskName: "New projects",
		TaskDetail: "You must lead the project and finish it",
		Date: "2022-01-01",
	}
	tasks = append(tasks, task)
	task1 := Tasks{
		ID: "2",
		TaskName: "Power project",
		TaskDetail: "We need to hire more staff",
		Date: "2022-01-01",
	}
	tasks = append(tasks, task1)
	fmt.Println("Your tasks are", tasks)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Home Page")
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	taskId := mux.Vars(r)
	fmt.Println(taskId["id"])
	var flag = false
	for i := 0; i < len(tasks); i ++ {
		if taskId["id"] == tasks[i].ID {
			json.NewEncoder(w).Encode(tasks[i])
			flag = true
			break
		}
	}
	if !flag {
		json.NewEncoder(w).Encode(map[string]string{"status":"Error"})
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task Tasks 
	// r.Body refers to the body in postData in the service.dart file
	// moves the body to the task variable
	// the variable is empty because it's being saving in the memory as a reference (&task)
	_ = json.NewDecoder(r.Body).Decode(&task)
	// creates a random ID
	task.ID = strconv.Itoa(rand.Intn(1000))
	currentTime := time.Now().Format("01-02-2006")
	task.Date = currentTime
	tasks = append(tasks, task)
	json.NewEncoder(w).Encode(task)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	flag := false
	// loop through all the tasks
	for i, v := range tasks {
		// if the param id matches a task id, update the task
		if v.ID == params["id"] {
			// takes everything up to the matching task (excluding the matching task) -> tasks[:i]
			// then gets the rest of the tasks (because of the spread operator ...) -> tasks[i+1:]...
			tasks = append(tasks[:i], tasks[i+1:]...)
			flag = true
			json.NewEncoder(w).Encode(map[string]string{"status":"Success"})
			return
		}
	}
	if !flag {
		json.NewEncoder(w).Encode(map[string]string{"status":"Error"})
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	flag := false
	// loop through all the tasks
	for i, v := range tasks {
		// if the param id matches a task id, update the task
		if v.ID == params["id"] {
			// takes everything up to the matching task (excluding the matching task) -> tasks[:i]
			// then gets the rest of the tasks (because of the spread operator ...) -> tasks[i+1:]...
			// so if we're updating the third index, it will grab: 0, 1, 2, 4, 5, etc...
			tasks = append(tasks[:i], tasks[i+1:]...)
			var task Tasks
			_ = json.NewDecoder(r.Body).Decode(&task)
			// sets the id because it deletes it in the previous append method
			task.ID = params["id"]
			currentTime := time.Now().Format("01-02-2006")
			task.Date = currentTime
			tasks = append(tasks, task)
			flag = true
			json.NewEncoder(w).Encode(task)
			return
		}
	}
	if !flag {
		json.NewEncoder(w).Encode(map[string]string{"status":"Error"})
	}
}

func handleRoutes() {
	router := mux.NewRouter()
	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/gettasks", getTasks).Methods("GET")
	router.HandleFunc("/gettask/", getTask).Queries("id" ,"{id}").Methods("GET")
	router.HandleFunc("/create", createTask).Methods("POST")
	router.HandleFunc("/delete/", deleteTask).Queries("id" ,"{id}").Methods("DELETE")
	router.HandleFunc("/update/", updateTask).Queries("id" ,"{id}").Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	allTasks()
	// fmt.Println("Hello, world!")
	handleRoutes()
}