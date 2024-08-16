package main

import (
	"html/template"
	"net/http"
	"strings"
)

var tmpl *template.Template

func init() {
	tmpl, _ = template.ParseGlob("templates/*.html")
}

func main() {
	initDB()
	defer db.Close()

	gRouter := InitRoutes()
	http.ListenAndServe(":8000", gRouter)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

func convertDoneToBool(isDone string) bool {
	var taskIsDone bool

	switch strings.ToLower(isDone) {
	case "yes", "on":
		taskIsDone = true
	case "no", "off":
		taskIsDone = false
	default:
		taskIsDone = false
	}

	return taskIsDone
}
