package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name"`
	SphereID primitive.ObjectID `bson:"sphere_id" json:"sphere_id"`
	ParentID primitive.ObjectID `bson:"parent_id, omitempty" json:"parent_id"`
}
