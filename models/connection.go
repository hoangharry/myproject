package models

import(
	"fmt"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitiateMongoClient() *mongo.Client {
	var err error
	var client *mongo.Client
	uri := "mongodb://HMH:bachkhoamt1@cluster0-shard-00-00.gjysk.mongodb.net:27017,cluster0-shard-00-01.gjysk.mongodb.net:27017,cluster0-shard-00-02.gjysk.mongodb.net:27017/<dbname>?ssl=true&replicaSet=atlas-8sqsyx-shard-0&authSource=admin&retryWrites=true&w=majority"
	opts := options.Client()
	opts.ApplyURI(uri)
	opts.SetMaxPoolSize(10)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if client, err = mongo.Connect(ctx, opts); err != nil {
		fmt.Println(err.Error())

	}
	return client
}