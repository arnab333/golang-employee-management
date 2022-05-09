package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoCollections struct {
	roles string
	user  string
}

var collectionNames = mongoCollections{
	roles: "roles",
	user:  "users",
}

var mongoOnce sync.Once

var redisOnce sync.Once

type MongoDBConnection struct {
	Database *mongo.Database
}

var DBConn MongoDBConnection

type redisConnection struct {
	redisClient *redis.Client
}

var redisConn redisConnection

func InitMongoConnection() func() {
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
		closeMongoConnection(client, ctx, cancel)
	}
}

func closeMongoConnection(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer func() {
		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal("Disconnect Error ==>", err)
		}
	}()

	defer cancel()
}

func InitRedisConnection() func() {
	redisOnce.Do(func() {
		dsn := os.Getenv("REDIS_DSN")

		redisConn.redisClient = redis.NewClient(&redis.Options{
			Addr:     dsn,                         //redis port
			Password: os.Getenv("REDIS_PASSWORD"), // no password set
			DB:       0,                           // use default DB
		})
		result, err := redisConn.redisClient.Ping(context.Background()).Result()
		if err != nil {
			panic(err)
		}
		log.Println("redis ==>", result)
	})

	return func() {
		closeRedisConnection(redisConn.redisClient)
	}
}

func closeRedisConnection(client *redis.Client) {
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatal("Close Error ==>", err)
		}
	}()
}
