package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Note struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Title        string               `bson:"title" json:"title"`
	CategoryIDs  []primitive.ObjectID `bson:"category_ids" json:"category_ids"`
	SphereID     primitive.ObjectID   `bson:"sphere_id" json:"sphere_id"`
	Content      string               `bson:"content" json:"content"`
	ExternalLink string               `bson:"external_link" json:"external_link"`
	InternalLink string               `bson:"internal_link" json:"internal_link"`
}
