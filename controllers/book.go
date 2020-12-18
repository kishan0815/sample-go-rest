package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sample-go-rest/config"
	"sample-go-rest/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookController struct {
	bookCollection *mongo.Collection
}

func NewBookController(bookCollection *mongo.Collection) *BookController {
	return &BookController{bookCollection}
}

// Get all books
func (bc BookController) GetBooks(rw http.ResponseWriter, req *http.Request) {
	var results []*models.Book
	results, err := getAll(bc)
	if err != nil {
		http.Error(rw, "Error fetching books", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(results)
}

func getAll(bc BookController) ([]*models.Book, error) {
	// passing bson.D{{}} matches all documents in the collection
	filter := bson.D{{}}
	return filterBooks(filter, bc)
}

func filterBooks(filter interface{}, bc BookController) ([]*models.Book, error) {
	// A slice of books
	var books []*models.Book

	cur, err := bc.bookCollection.Find(config.AppContext, filter)
	if err != nil {
		return books, err
	}

	for cur.Next(config.AppContext) {
		var b models.Book
		err := cur.Decode(&b)
		if err != nil {
			return books, err
		}

		books = append(books, &b)
	}

	if err := cur.Err(); err != nil {
		return books, err
	}

	// once exhausted, close the cursor
	cur.Close(config.AppContext)

	if len(books) == 0 {
		return books, mongo.ErrNoDocuments
	}

	return books, nil
}

// Get a book
func (bc BookController) GetBook(rw http.ResponseWriter, req *http.Request) {
	var book models.Book
	rw.Header().Set("Content-Type", "application/json")
	id := mux.Vars(req)["id"]
	objID, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": objID}
	err := bc.bookCollection.FindOne(config.AppContext, filter).Decode(&book)
	if err != nil {
		fmt.Println(err)
		http.Error(rw, "Error fetching book", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(rw).Encode(book)
}

// Create a book
func (bc BookController) CreateBook(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	var book models.Book
	json.NewDecoder(req.Body).Decode(&book)
	book.ID = primitive.NewObjectID()
	id, err := bc.bookCollection.InsertOne(config.AppContext, book)
	if err != nil {
		http.Error(rw, "Error creating book", http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	fmt.Println(id.InsertedID)
	json.NewEncoder(rw).Encode(book)
}

// Update a existing book
func (bc BookController) UpdateBook(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	id := mux.Vars(req)["id"]
	objID, _ := primitive.ObjectIDFromHex(id)
	var book models.Book
	json.NewDecoder(req.Body).Decode(&book)
	book.ID = objID
	_, err := bc.bookCollection.ReplaceOne(
		config.AppContext,
		bson.M{"_id": objID},
		book,
	)
	if err != nil {
		fmt.Println(err)
	}
	books, _ := getAll(bc)
	json.NewEncoder(rw).Encode(books)
}

// Delete a book
func (bc BookController) DeleteBook(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	id := mux.Vars(req)["id"]
	objID, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": objID}

	res, err := bc.bookCollection.DeleteOne(config.AppContext, filter)
	if err != nil {
		http.Error(rw, "Error deleting book", http.StatusInternalServerError)
		return
	}

	if res.DeletedCount == 0 {
		http.Error(rw, "No tasks were deleted", http.StatusInternalServerError)
		return
	}

	books, _ := getAll(bc)
	json.NewEncoder(rw).Encode(books)
}
