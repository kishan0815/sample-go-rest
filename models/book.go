package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Book Struct
type Book struct {
	ID     primitive.ObjectID `bson:"_id"`
	Title  string             `bson:"title"`
	Author *Author            `bson:"author"`
}
