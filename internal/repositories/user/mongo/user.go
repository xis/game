package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type userRecord struct {
	ID           primitive.ObjectID `bson:"_id"`
	Username     string             `bson:"username"`
	PasswordHash string             `bson:"passwordHash"`
}
