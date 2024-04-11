package models

import "go.mongodb.org/mongo-driver/x/mongo/driver/uuid"

type Question struct {
	ID       uuid.UUID
	Title    string
	Category Category
	Theme    Theme
	Content  string
}

type Theme struct {
	ID    uuid.UUID
	Title string
}

type Category struct {
	ID    uuid.UUID
	Title string
}
