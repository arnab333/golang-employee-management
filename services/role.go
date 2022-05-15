package services

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRole struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Permissions []string           `bson:"permissions" json:"permissions"`
}

func (conn *MongoDBConnection) FindRole(c *gin.Context, filters interface{}) (UserRole, error) {
	collection := conn.Database.Collection(collectionNames.roles)

	var data UserRole

	if filters == nil {
		filters = bson.M{}
	}

	err := collection.FindOne(c, filters).Decode(&data)

	return data, err
}

func (conn *MongoDBConnection) FindRoles(c *gin.Context, filters interface{}) ([]UserRole, error) {
	collection := conn.Database.Collection(collectionNames.roles)

	var data []UserRole

	if filters == nil {
		filters = bson.M{}
	}

	cur, err := collection.Find(c, filters)
	if err != nil {
		return nil, err
	}

	err = cur.All(c, &data)

	if err != nil {
		return nil, err
	}

	return data, err
}

func (conn *MongoDBConnection) UpdateRole(c *gin.Context, filters interface{}, data interface{}) (*mongo.UpdateResult, error) {
	collection := conn.Database.Collection(collectionNames.roles)

	opts := options.Update().SetUpsert(true)

	if filters == nil {
		filters = bson.M{}
	}

	return collection.UpdateOne(c, filters, data, opts)
}
