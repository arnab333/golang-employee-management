package services

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRole struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name" json:"name" binding:"required"`
}

func (conn *MongoDBConnection) InsertUserRole(data interface{}) (*mongo.InsertOneResult, error) {
	collection := conn.Database.Collection(collectionNames.userRole)

	return collection.InsertOne(conn.Context, data)
}
