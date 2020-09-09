package models

import(
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"time"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	// "encoding/binary"
	// "labix.org/v2/mgo/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo/gridfs"
	"log"
	// "bytes"
	// "io/ioutil"
	// "crypto/md5"
	// "encoding/json"

)


type Image struct {
	ID  int `json:"id" bson:"_id"`
	Direction string `json:"direction" bson:"direction"`
	Time string `json:"time" bson:"time"`
	FileId primitive.ObjectID `bson:"evidence" json:"evidence"`
	Data primitive.Binary `bson:"data" json:"data"`
}

type ImgChunk struct {
	Data primitive.Binary `bson:"data" json:"data"`
}

type ImageManager struct {
	imgs []*Image
	client *mongo.Client
}

var DefaultImageList *ImageManager

func NewImageManager(client *mongo.Client) *ImageManager {
	return &ImageManager{client:client}
}



func (m *ImageManager) All() []*Image {
	db := m.client.Database("test")
	collection := db.Collection("test")
	// fileCol := db.Collection("fs.files")
	chunkCol := db.Collection("fs.chunks")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}	
	defer cursor.Close(ctx)
	var tmp []*Image
	for cursor.Next(ctx){

		// var news bson.M
		// cursor.Decode(&news)
		// newsByte , _ := bson.Marshal(news)
		var newimg *Image
		var newImg Image 
		err = cursor.Decode(&newImg)
		if err != nil {
			log.Fatal(err)
		}
		// var imgFile ImgFile
		var imgChunk ImgChunk
		// err = fileCol.FindOne(ctx, bson.M{"_id": newImg.FileId}).Decode(&imgFile)
		if err != nil {
			log.Fatal(err)
		}
		filter := &bson.M{"files_id": newImg.FileId}
		err = chunkCol.FindOne(ctx, filter).Decode(&imgChunk)

		newImg.Data = imgChunk.Data
		newimg = &newImg
		tmp = append(tmp, newimg)
			

	}
	if (len(tmp) > len(m.imgs)){
		m.imgs = tmp
	}

	return m.imgs
}

func init(){
	uri := "mongodb://HMH:bachkhoamt1@cluster0-shard-00-00.gjysk.mongodb.net:27017,cluster0-shard-00-01.gjysk.mongodb.net:27017,cluster0-shard-00-02.gjysk.mongodb.net:27017/<dbname>?ssl=true&replicaSet=atlas-8sqsyx-shard-0&authSource=admin&retryWrites=true&w=majority"
	// uri := "mongodb://localhost:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successful connected and pinged")
	DefaultImageList = NewImageManager(client)
	
}
