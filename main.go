package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Task struct {
	Id   int
	Task string
	Done bool
}

var db *sql.DB
var tmpl *template.Template

func init() {
	tmpl, _ = template.ParseGlob("templates/*.html")
}

func initDB() {
	var err error

	errENV := godotenv.Load(".env")
	if errENV != nil {
		log.Fatal("Error with .env loading", errENV)
	}
	db_connection_string := os.Getenv("DB_CONNECTION_STRING")

	db, err = sql.Open("postgres", db_connection_string)

	if err != nil {
		log.Fatal("Error Connecting to DB start", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Error Connecting to DB end", err)
	}

	fmt.Println("db connected succes")
}

func main() {

	gRouter := mux.NewRouter()

	initDB()
	defer db.Close()

	gRouter.HandleFunc("/", homeHandler)

	gRouter.HandleFunc("/tasks", fetchTasks).Methods("GET")

	gRouter.HandleFunc("/tasks", addNewTask).Methods("POST")

	gRouter.HandleFunc("/getTaskForm", getTaskForm)

	http.ListenAndServe(":8000", gRouter)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

func fetchTasks(w http.ResponseWriter, r *http.Request) {
	todos, errDB := getTasks()
	if errDB != nil {
		log.Fatal("Error fetching tasks: ", errDB)
	}

	err := tmpl.ExecuteTemplate(w, "todoList", todos)

	if err != nil {
		panic(err)
	}
}

func getTasks() ([]Task, error) {
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

func getTaskForm(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "addTaskForm", nil)
}

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

	todos, _ := getTasks()
	tmpl.ExecuteTemplate(w, "todoList", todos)
}
