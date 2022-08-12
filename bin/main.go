package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"main.go/service"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/posts", service.GetAllPost).Methods("GET")
	router.HandleFunc("/posts", service.InsertsPost).Methods("POST")
	router.HandleFunc("/posts", service.DeletePost).Methods("DELETE")
	router.HandleFunc("/posts", service.UpdatePost).Methods("PUT")
	router.HandleFunc("/users", service.GetAllUser).Methods("GET")
	router.HandleFunc("/users", service.InsertUser).Methods("POST")
	router.HandleFunc("", Initial).Methods("GET")

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))

}

func Initial(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Connection Established")
}
