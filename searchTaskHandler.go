package main

import (
	"net/http"
	"strings"
)

func searchTask(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("searchQuery")
	lowercaseQuery := strings.ToLower(query)

	rows, err := db.Query("SELECT id, task, done FROM tasks WHERE LOWER(task) LIKE '%' || $1 || '%'", lowercaseQuery)
	if err != nil {
		http.Error(w, "Failed to find task", http.StatusInternalServerError)
	}

	defer rows.Close()

	var todos []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Task, &task.Done); err != nil {
			http.Error(w, "Failed to scan task", http.StatusInternalServerError)
			return
		}
		todos = append(todos, task)
	}

	err = tmpl.ExecuteTemplate(w, "todoList", todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
