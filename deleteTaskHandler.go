package main

import (
	"net/http"
)

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/deleteTask/"):]

	_, err := db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
