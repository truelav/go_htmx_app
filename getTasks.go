package main

import (
	"log"
	"net/http"
)

func fetchTasks(w http.ResponseWriter, r *http.Request) {
	todos, errDB := getTasksDB()
	if errDB != nil {
		log.Fatal("Error fetching tasks: ", errDB)
	}

	err := tmpl.ExecuteTemplate(w, "todoList", todos)

	if err != nil {
		panic(err)
	}
	// fmt.Println(todos)
}

func getTasksDB() ([]Task, error) {
	query := "SELECT id, task, done FROM tasks"

	rows, errDB := db.Query(query)

	if errDB != nil {
		log.Fatal("Error querying DB: ", errDB)
	}

	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task
		rowErr := rows.Scan(&task.Id, &task.Task, &task.Done)

		if rowErr != nil {
			log.Fatal("Error something wrong scanning rows: ", errDB)
		}

		tasks = append(tasks, task)
	}

	return tasks, errDB
}

func getTaskDB(taskId string) (Task, error) {

	var task Task
	query := "SELECT id, task, done FROM tasks WHERE id = $1"
	err := db.QueryRow(query, taskId).Scan(&task.Id, &task.Task, &task.Done)

	return task, err
}

func getEditTaskForm(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Path[len("/getEditTaskForm/"):]

	task, err := getTaskDB(taskId)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
	}

	tmpl.ExecuteTemplate(w, "editTaskForm", task)
}

func getTaskForm(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "addTaskForm", nil)
}
