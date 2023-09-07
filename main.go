package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	handlers "go-api/handlers"
)

var (
	DATABASE_URL, DB_DRIVER, PORT string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Couldn't load env on startup!!")
	}
	DATABASE_URL = os.Getenv("DATABASE_URL")
	DB_DRIVER = os.Getenv("DB_DRIVER")
	PORT = os.Getenv("PORT")
}

func DBClient() (*sql.DB, error) {
	db, err := sql.Open(DB_DRIVER, DATABASE_URL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Connected to DB")
	return db, nil
}

type Server struct {
	Router *chi.Mux
	DB     *sql.DB
}

func CreateServer(db *sql.DB) *Server {
	server := &Server{
		Router: chi.NewRouter(),
		DB:     db,
	}
	return server
}

func Greet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!!"))
}

func (server *Server) MountHandlers() {
	server.Router.Get("/greet", Greet)

	todosRouter := chi.NewRouter()
	todosRouter.Group(func(r chi.Router) {
		r.Get("/", server.GetTodos)
		r.Post("/", server.AddTodo)
	})

	userRouter := chi.NewRouter()
	userRouter.Post("/login", handlers.LoginUser)
	userRouter.Post("/signup", handlers.CreateUser)
	userRouter.Group(func(r chi.Router) {
		r.Get("/{id}", handlers.GetUser)
		r.Post("/{id}", handlers.UpdateUser)
		r.Delete("/{id}", handlers.DeleteUser)
	})

	server.Router.Mount("/todos", todosRouter)
	server.Router.Mount("/user", userRouter)
}

func main() {
	db, err := DBClient()
	if err != nil {
		log.Fatalln(err)
	}
	server := CreateServer(db)
	server.MountHandlers()

	fmt.Println("server running on port:5000")
	http.ListenAndServe(PORT, server.Router)
}

type Todo struct {
	Id        int       `json:"id"`
	Task      string    `json:"task"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TodoRequestBody struct {
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

func scanRow(rows *sql.Rows) (*Todo, error) {
	todo := new(Todo)
	err := rows.Scan(&todo.Id,
		&todo.Task,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (server *Server) AddTodo(w http.ResponseWriter, r *http.Request) {
	todo := new(TodoRequestBody)
	if err := json.NewDecoder(r.Body).Decode(todo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please enter a correct Todo!!"))
		return
	}

	query := `INSERT INTO Todos (task, completed) VALUES (?, ?)`
	_, err := server.DB.Exec(query, todo.Task, todo.Completed)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something bad happened on the server :("))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Todo added!!"))
}

func (server *Server) GetTodos(w http.ResponseWriter, r *http.Request) {
	query := `SELECT * FROM Todos ORDER BY created_at DESC`

	rows, err := server.DB.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something bad happened on the server :("))
		return
	}

	var todos []*Todo

	for rows.Next() {
		todo, err := scanRow(rows)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something bad happened on the server :("))
			return
		}
		todos = append(todos, todo)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}
