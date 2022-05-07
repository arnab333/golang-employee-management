package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoCollections struct {
	userRole string
	user     string
}

var collectionNames = mongoCollections{
	userRole: "roles",
	user:     "users",
}

var mongoOnce sync.Once

type MongoDBConnection struct {
	Database *mongo.Database
}

var DBConn MongoDBConnection

func GetMongoConnection() func() {
	var client *mongo.Client
	var cancel context.CancelFunc
	var ctx context.Context

	mongoOnce.Do(func() {
		username := os.Getenv("MONGO_USERNAME")
		pass := os.Getenv("MONGO_PASSWORD")
		dbname := os.Getenv("MONGO_DBNAME")

		if username == "" || pass == "" || dbname == "" {
			log.Fatal("You must set your 'MONGO_USERNAME' environmental variable.")
		}

		client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@cluster0.6x1pj.mongodb.net/%s?retryWrites=true&w=majority", os.Getenv("MONGO_USERNAME"), os.Getenv("MONGO_PASSWORD"), os.Getenv("MONGO_DBNAME"))))
		if err != nil {
			log.Fatal("NewClient Error ==>", err)
		}

		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
		// defer cancel()
		if err = client.Connect(ctx); err != nil {
			log.Fatal("Connect Error ==>", err)
		}

		// defer func() {
		// 	// client.Disconnect method also has deadline.
		// 	// returns error if any,
		// 	if err := client.Disconnect(ctx); err != nil {
		// 		log.Fatal("Disconnect Error ==>", err)
		// 	}
		// }()

		DBConn.Database = client.Database(dbname)

		// ## The following part is only for checking if the database connection is successfull or not.
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Fatal("Ping Error ==>", err)
		}

		fmt.Print("Ping Succes!!")
	})

	return func() {
		CloseMongoConnection(client, ctx, cancel)
	}
}

func CloseMongoConnection(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer func() {
		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal("Disconnect Error ==>", err)
		}
	}()

	defer cancel()
}
