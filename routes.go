package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	gRouter := mux.NewRouter()

	gRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	gRouter.HandleFunc("/", homeHandler).Methods("GET")

	gRouter.HandleFunc("/tasks", fetchTasks).Methods("GET")

	gRouter.HandleFunc("/tasks", addNewTask).Methods("POST")

	gRouter.HandleFunc("/editTask/{{.Id}}", editTask)

	gRouter.HandleFunc("/toggleTask/{{.Id}}", toggleTaskDone)

	gRouter.HandleFunc("/deleteTask/{{.Id}}", deleteTask)

	gRouter.HandleFunc("/cancelEditTask/{{.Id}}", cancelEditTask)

	gRouter.HandleFunc("/searchTask", searchTask)

	gRouter.HandleFunc("/getEditTaskForm/{{Id}}", getEditTaskForm)

	gRouter.HandleFunc("/getTaskForm", getTaskForm)

	return gRouter
}
