package services

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (conn *MongoDBConnection) UpdateUser(c *gin.Context, filters interface{}, data interface{}) (*mongo.UpdateResult, error) {
	collection := conn.Database.Collection(collectionNames.users)

	opts := options.Update()

	if filters == nil {
		filters = bson.M{}
	}

	return collection.UpdateOne(c, filters, data, opts)
}

func (conn *MongoDBConnection) FindUsers(c *gin.Context, filters interface{}, limit, pageNo int64) ([]User, error) {
	var opts *options.FindOptions

	collection := conn.Database.Collection(collectionNames.users)

	var data []User

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
		var el User
		if err := cur.Decode(&el); err != nil {
			log.Println(err)
		}

		data = append(data, el)
	}

	return data, nil

}
