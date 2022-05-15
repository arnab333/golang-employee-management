package services

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName string             `bson:"fullName" json:"fullName" binding:"required"`
	Email    string             `bson:"email" json:"email" binding:"required"`
	Password string             `bson:"password" json:"password" binding:"required"`
	Role     string             `bson:"role" json:"role" binding:"required"`
	IsActive bool               `bson:"isActive" json:"isActive"`
}

func (conn *MongoDBConnection) InsertUser(c *gin.Context, data interface{}) (*mongo.InsertOneResult, error) {
	collection := conn.Database.Collection(collectionNames.users)

	return collection.InsertOne(c, data)
}

func (conn *MongoDBConnection) FindUser(c *gin.Context, filters interface{}) (User, error) {
	collection := conn.Database.Collection(collectionNames.users)

	var data User

	if filters == nil {
		filters = bson.M{}
	}

	err := collection.FindOne(c, filters).Decode(&data)

	return data, err
}
