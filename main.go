package main

import (
	"html/template"
	"net/http"
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
