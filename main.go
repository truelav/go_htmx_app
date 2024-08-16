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

	gRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	gRouter.HandleFunc("/", homeHandler)

	gRouter.HandleFunc("/tasks", fetchTasks).Methods("GET")

	gRouter.HandleFunc("/tasks", addNewTask).Methods("POST")

	gRouter.HandleFunc("/editTask/{{.Id}}", editTask)

	gRouter.HandleFunc("/toggleTask/{{.Id}}", toggleTaskDone)

	gRouter.HandleFunc("/deleteTask/{{.Id}}", deleteTask)

	gRouter.HandleFunc("/cancelEditTask/{{.Id}}", cancelEditTask)

	gRouter.HandleFunc("/searchTask", searchTask)

	gRouter.HandleFunc("/getEditTaskForm/{{Id}}", getEditTaskForm)

	gRouter.HandleFunc("/getTaskForm", getTaskForm)

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
	// fmt.Println(todos)
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

func getTaskDB(taskId string) (Task, error) {

	var task Task
	query := "SELECT id, task, done FROM tasks WHERE id = $1"
	err := db.QueryRow(query, taskId).Scan(&task.Id, &task.Task, &task.Done)

	return task, err
}

func getEditTaskForm(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Path[len("/getEditTaskForm/"):]

	task, err := getTaskDB(taskId)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
	}

	tmpl.ExecuteTemplate(w, "editTaskForm", task)
}

func cancelEditTask(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Path[len("/cancelEditTask/"):]

	task, err := getTaskDB(taskId)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
	}

	tmpl.ExecuteTemplate(w, "task", task)
}

func toggleTaskDone(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/toggleTask/"):]
	taskIsDone := r.FormValue("done")
	done := convertDoneToBool(taskIsDone)

	_, err := db.Exec("UPDATE tasks SET done = $1 WHERE id = $2", done, id)
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

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/deleteTask/"):]

	_, err := db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	// Return an empty response to remove the task from the UI
	w.WriteHeader(http.StatusOK)
}

func searchTask(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("searchQuery")
	lowercaseQuery := strings.ToLower(query)
	// fmt.Println(lowercaseQuery)

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
