package models

// Author Struct
type Author struct {
	Firstname string `bson:"firstName"`
	LastName  string `bson:"lastName"`
}
