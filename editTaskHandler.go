package main

import (
	"net/http"
)

func editTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/editTask/"):]
	task := r.FormValue("task")
	taskIsDone := r.FormValue("done")

	done := convertDoneToBool(taskIsDone)

	_, err := db.Exec("UPDATE tasks SET task = $1, done = $2 WHERE id = $3", task, done, id)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	var updatedTask Task
	err = db.QueryRow("SELECT id, task, done FROM tasks WHERE id = $1", id).Scan(&updatedTask.Id, &updatedTask.Task, &updatedTask.Done)
	if err != nil {
		http.Error(w, "Failed to fetch updated task", http.StatusInternalServerError)
		return
	}

	tmpl.ExecuteTemplate(w, "task", updatedTask)
}
