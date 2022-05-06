package services

import (
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRole struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Value int64              `bson:"value"`
	Text  string             `bson:"text"`
}

func (conn *MongoDBConnection) InsertUserRole(data []interface{}) *mongo.InsertManyResult {
	collection := conn.Database.Collection(collectionNames.userRole)

	result, err := collection.InsertMany(conn.Context, data)

	if err != nil {
		log.Fatal("InsertMany Error ==>", err)
	}
	return result
}
