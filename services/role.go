package services

import (
	"log"

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

func (conn *MongoDBConnection) FindRoles(c *gin.Context, filters interface{}, limit, pageNo int64) ([]UserRole, error) {
	var opts *options.FindOptions

	collection := conn.Database.Collection(collectionNames.roles)

	var data []UserRole

	if filters == nil {
		filters = bson.M{}
	}

	if limit != 0 && pageNo != 0 {
		opts = newMongoPaginate(limit, pageNo).getPaginatedOpts()
	}

	cur, err := collection.Find(c, filters, opts)
	if err != nil {
		return nil, err
	}

	if opts == nil {
		err = cur.All(c, &data)

		if err != nil {
			return nil, err
		}

		return data, nil
	}

	for cur.Next(c) {
		var el UserRole
		if err := cur.Decode(&el); err != nil {
			log.Println(err)
		}

		data = append(data, el)
	}

	return data, nil
}

func (conn *MongoDBConnection) InsertRole(ctx *gin.Context, data interface{}) (*mongo.InsertOneResult, error) {
	collection := conn.Database.Collection(collectionNames.roles)

	opts := options.InsertOne()

	return collection.InsertOne(ctx, data, opts)
}

func (conn *MongoDBConnection) UpdateRole(c *gin.Context, filters interface{}, data interface{}) (*mongo.UpdateResult, error) {
	collection := conn.Database.Collection(collectionNames.roles)

	opts := options.Update()

	if filters == nil {
		filters = bson.M{}
	}

	return collection.UpdateOne(c, filters, data, opts)
}
