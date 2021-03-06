package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/arnab333/golang-employee-management/helpers"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoCollections struct {
	roles    string
	users    string
	holidays string
}

var collectionNames = mongoCollections{
	roles:    "roles",
	users:    "users",
	holidays: "holidays",
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

type mongoPaginate struct {
	limit  int64
	pageNo int64
}

func InitMongoConnection() func() {
	var client *mongo.Client
	var cancel context.CancelFunc
	var ctx context.Context

	mongoOnce.Do(func() {
		username := os.Getenv(helpers.EnvKeys.MONGO_USERNAME)
		pass := os.Getenv(helpers.EnvKeys.MONGO_PASSWORD)
		dbname := os.Getenv(helpers.EnvKeys.MONGO_DBNAME)

		if username == "" || pass == "" || dbname == "" {
			log.Fatal("You must set your 'MONGO_USERNAME' environmental variable.")
		}

		client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@cluster0.6x1pj.mongodb.net/%s?retryWrites=true&w=majority", username, pass, dbname)))
		if err != nil {
			log.Fatal("NewClient Error ==>", err)
		}

		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
		// defer cancel()
		if err = client.Connect(ctx); err != nil {
			log.Fatal("Connect Error ==>", err)
		}

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
		dsn := os.Getenv(helpers.EnvKeys.REDIS_DSN)

		redisConn.redisClient = redis.NewClient(&redis.Options{
			Addr:     dsn,                                       //redis port
			Password: os.Getenv(helpers.EnvKeys.REDIS_PASSWORD), // no password set
			DB:       0,                                         // use default DB
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

func newMongoPaginate(limit, pageNo int64) *mongoPaginate {
	return &mongoPaginate{
		limit:  limit,
		pageNo: pageNo,
	}
}

func (mp *mongoPaginate) getPaginatedOpts() *options.FindOptions {
	limit := mp.limit
	skip := (mp.pageNo * mp.limit) - mp.limit
	opts := options.FindOptions{Limit: &limit, Skip: &skip}

	return &opts
}
