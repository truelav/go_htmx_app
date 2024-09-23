package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// var tmpl *template.Template

// func init() {
// 	tmpl, _ = template.ParseGlob("templates/*.html")
// }

type Task struct {
	ID   int
	Name string
}

func main() {
	initDB()
	defer db.Close()

	// gRouter := InitRoutes()
	gRouter := mux.NewRouter()

	gRouter.HandleFunc("/", homeHandler).Methods("GET")
	gRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":8000", gRouter)
}

// func homeHandler(w http.ResponseWriter, r *http.Request) {
// 	tmpl.ExecuteTemplate(w, "index.html", nil)
// }

func homeHandler(w http.ResponseWriter, r *http.Request) {
	views.Home().Render(r.Context(), w)
}
