package services

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Holiday struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Summary  string             `bson:"summary" json:"summary" binding:"required"`
	Date     string             `bson:"date" json:"date" binding:"required"`
	IsActive bool               `bson:"isActive" json:"isActive"`
}

func (conn *MongoDBConnection) FindHoliday(ctx context.Context, filters interface{}) (Holiday, error) {
	collection := conn.Database.Collection(collectionNames.holidays)

	if filters == nil {
		filters = bson.M{}
	}

	var data Holiday

	err := collection.FindOne(ctx, filters).Decode(&data)

	return data, err
}

func (conn *MongoDBConnection) InsertHoliday(ctx context.Context, data interface{}) (*mongo.InsertOneResult, error) {
	collection := conn.Database.Collection(collectionNames.holidays)

	opts := options.InsertOne()

	return collection.InsertOne(ctx, data, opts)
}

func (conn *MongoDBConnection) FindHolidays(c *gin.Context, filters interface{}, limit, pageNo int64) ([]Holiday, error) {
	var opts *options.FindOptions

	collection := conn.Database.Collection(collectionNames.holidays)

	var data []Holiday

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
		var el Holiday
		if err := cur.Decode(&el); err != nil {
			log.Println(err)
		}

		data = append(data, el)
	}

	return data, nil
}
