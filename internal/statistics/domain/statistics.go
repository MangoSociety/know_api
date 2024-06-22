package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Statistics struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Status    int                `bson:"status"`
	NoteID    primitive.ObjectID `bson:"note_id"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}
