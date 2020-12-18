package main

import (
	"log"
	"net/http"

	"sample-go-rest/config"
	"sample-go-rest/controllers"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// Init books var as a slice of Book struct
//var books []Book

var collection *mongo.Collection

func init() {
	collection = config.BookCollection
}

func main() {
	// Init Router
	r := mux.NewRouter()
	bc := controllers.NewBookController(collection)
	// Route Handlers
	r.HandleFunc("/api/books", bc.GetBooks).Methods("GET")
	r.HandleFunc("/api/book/{id}", bc.GetBook).Methods("GET")
	r.HandleFunc("/api/book", bc.CreateBook).Methods("POST")
	r.HandleFunc("/api/book/{id}", bc.UpdateBook).Methods("PUT")
	r.HandleFunc("/api/book/{id}", bc.DeleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}
