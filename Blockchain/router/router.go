package router

import (
	"../server"
	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()
	router.HandleFunc("/api/block", server.GetAllBlock).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/block", server.CreateBlock).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/task", server.CreateTask).Methods("POST", "OPTIONS")
	//	router.HandleFunc("/api/task", server.CreateTask).Methods("POST", "OPTIONS")
	//	router.HandleFunc("/api/task/{id}", server.TaskComplete).Methods("PUT", "OPTIONS")
	//	router.HandleFunc("/api/undoTask/{id}", server.UndoTask).Methods("PUT", "OPTIONS")
	//	router.HandleFunc("/api/deleteTask/{id}", server.DeleteTask).Methods("DELETE", "OPTIONS")
	//router.HandleFunc("/api/deleteAllTask", server.DeleteAllTask).Methods("DELETE", "OPTIONS")
	return router
}
