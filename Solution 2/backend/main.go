package main

import (
	"backend/config"
	"backend/controllers/subtaskcontroller"
	"backend/controllers/taskcontroller"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// 	Database connection
	username := "sql12648296"
	password := "a9S1uYL7md"
	hostname := "sql12.freemysqlhosting.net"
	dbname := "sql12648296"

	config.ConnectDB(username, password, hostname, dbname)

	// Create a new CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Replace with your frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Create a new Gorilla Mux router
	router := mux.NewRouter()

	// Define your API routes here

	// 1. Task
	router.HandleFunc("/task/ongoing", taskcontroller.GetAllDataOnGoing)
	router.HandleFunc("/task/complete", taskcontroller.GetAllDataComplete)
	router.HandleFunc("/task/create", taskcontroller.CreateTask)
	router.HandleFunc("/task/update/{id}", taskcontroller.UpdateTaskNameByID).Methods("PUT")
	router.HandleFunc("/task/deadline/{id}", taskcontroller.UpdateDeadlineTaskById).Methods("PUT")
	router.HandleFunc("/task/delete/{id}", taskcontroller.DeleteTaskById).Methods("DELETE")
	router.HandleFunc("/task/complete/{id}", taskcontroller.UpdateTaskCompleteByID).Methods("PUT")

	// 2. Subtask
	router.HandleFunc("/subtask/create", subtaskcontroller.CreateSubTask)
	router.HandleFunc("/subtask/update/{id}", subtaskcontroller.UpdateSubTaskNameByID).Methods("PUT")
	router.HandleFunc("/subtask/delete/{id}", subtaskcontroller.DeleteSubTaskById).Methods("DELETE")
	router.HandleFunc("/subtask/complete/{id}", subtaskcontroller.UpdateSubTaskCompleteById).Methods("PUT")

	// Use the CORS middleware with your router
	handler := c.Handler(router)

	// Run server
	log.Println("Server running on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
