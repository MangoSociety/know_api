package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserSelection struct {
	ID                 primitive.ObjectID  `bson:"_id,omitempty"`
	ChatID             int                 `bson:"chat_id"`
	SelectedCategories []SelectedParameter `bson:"selected_parameters"`
}

type SelectedCategory struct {
	CategoryID    primitive.ObjectID `bson:"category_id"`
	CategoryTitle string             `bson:"category_title"`
}

type SelectedParameter struct {
	SphereID    primitive.ObjectID `bson:"sphere_id"`
	SphereTitle string             `bson:"sphere_title"`
	Categories  []SelectedCategory `bson:"categories"`
}
