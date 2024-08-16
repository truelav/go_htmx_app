package main

import (
	"net/http"
)

func cancelEditTask(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Path[len("/cancelEditTask/"):]

	task, err := getTaskDB(taskId)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
	}

	tmpl.ExecuteTemplate(w, "task", task)
}
