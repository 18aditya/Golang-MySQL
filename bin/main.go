package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"main.go/middleware"
	"main.go/service"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	router.Handle("/posts", middleware.ValidateJWT(service.GetAllPost)).Methods("GET")
	router.Handle("/posts", middleware.ValidateJWT(service.InsertsPost)).Methods("POST")
	router.Handle("/posts", middleware.ValidateJWT(service.DeletePost)).Methods("DELETE")
	router.Handle("/posts", middleware.ValidateJWT(service.UpdatePost)).Methods("PUT")
	router.Handle("/users", middleware.ValidateJWT(service.GetAllUser)).Methods("GET")
	router.Handle("/users", middleware.ValidateJWT(service.InsertUser)).Methods("POST")
	router.Handle("/jwt", middleware.ValidateJWT(service.GetJWT))
	router.HandleFunc("/", Initial).Methods("GET")

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))

}

func Initial(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Connection Established")

}
