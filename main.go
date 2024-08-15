package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

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

	gRouter.HandleFunc("/editTask", editTask)

	gRouter.HandleFunc("/getTaskForm", getTaskForm)

	gRouter.HandleFunc("/cancelEdit/{{.ID}}", cancelEditTask)

	gRouter.HandleFunc("/getEditTaskForm/{{Id}}", getEditTaskForm)

	http.ListenAndServe(":8000", gRouter)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

func fetchTasks(w http.ResponseWriter, r *http.Request) {
	todos, errDB := getTasksDB()
	if errDB != nil {
		log.Fatal("Error fetching tasks: ", errDB)
	}

	err := tmpl.ExecuteTemplate(w, "todoList", todos)

	if err != nil {
		panic(err)
	}
}

func getTaskDB(taskId string) (Task, error) {

	var task Task
	query := "SELECT id, task, done FROM tasks WHERE id = $1"
	err := db.QueryRow(query, taskId).Scan(&task.Id, &task.Task, &task.Done)

	return task, err
}

func getTasksDB() ([]Task, error) {
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

func getEditTaskForm(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Path[len("/getEditTaskForm/"):]

	task, err := getTaskDB(taskId)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
	}

	tmpl.ExecuteTemplate(w, "editTaskForm", task)
}

func getTaskItem(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Path[len("/cancelEditTask/"):]

	task, err := getTaskDB(taskId)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
	}

	tmpl.ExecuteTemplate(w, "task", task)
}

func editTask(w http.ResponseWriter, r *http.Request) {
	task := r.FormValue("task")
	taskId := r.FormValue("id")
	done := r.FormValue("done")
	taskIsDone := convertDoneToBool(done)

	fmt.Println(task, taskId, taskIsDone)
}

func cancelEditTask(w http.ResponseWriter, r *http.Request) {

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

	todos, _ := getTasksDB()
	tmpl.ExecuteTemplate(w, "todoList", todos)
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
