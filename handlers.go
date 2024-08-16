package main

import (
	"fmt"
	"net/http"
)

func toggleTaskDone(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/toggleTask/"):]
	taskIsDone := r.FormValue("done")
	done := convertDoneToBool(taskIsDone)

	_, err := getDB().Exec("UPDATE tasks SET done = $1 WHERE id = $2", done, id)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	var updatedTask Task
	err = db.QueryRow("SELECT id, task, done  FROM tasks WHERE id = $1", id).Scan(&updatedTask.Id, &updatedTask.Task, &updatedTask.Done)
	if err != nil {
		http.Error(w, "Failed to fetch updated task", http.StatusInternalServerError)
		return
	}

	fmt.Println(updatedTask)
	tmpl.ExecuteTemplate(w, "task", updatedTask)
}
