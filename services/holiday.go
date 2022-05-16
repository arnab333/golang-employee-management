package services

import (
	"context"

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

func (conn *MongoDBConnection) FindHolidays(c *gin.Context, filters interface{}) ([]Holiday, error) {
	collection := conn.Database.Collection(collectionNames.holidays)

	var data []Holiday

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
