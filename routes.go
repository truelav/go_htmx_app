package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
)

func createServer() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func InitRoutes() *mux.Router {
	gRouter := mux.NewRouter()

	gRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	gRouter.HandleFunc("/", homeHandler)

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
