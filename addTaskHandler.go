package main

import (
	"net/http"
)

func addNewTask(w http.ResponseWriter, r *http.Request) {
	task := r.FormValue("task")

	if task == "" {
		http.Error(w, "Task cannot be empty", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO tasks (task, done) VALUES ($1, $2)"

	_, err := db.Exec(query, task, false)
	if err != nil {
		http.Error(w, "Failed to add task", http.StatusInternalServerError)
		return
	}

	todos, _ := getTasksDB()
	tmpl.ExecuteTemplate(w, "todoList", todos)
}
