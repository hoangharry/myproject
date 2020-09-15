package models

import(
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"context"
	"time"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)


type Image struct {
	ID  int `json:"id" bson:"_id"`
	Direction string `json:"direction" bson:"direction"`
	Time string `json:"time" bson:"time"`
	Line int `json:"line" bson:"line"`
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
	chunkCol := db.Collection("fs.chunks")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}	
	defer cursor.Close(ctx)
	var tmp []*Image
	for cursor.Next(ctx){
		var newimg *Image
		var newImg Image 
		err = cursor.Decode(&newImg)
		if err != nil {
			fmt.Println(err)
		}
		var imgChunk ImgChunk
		filter := &bson.M{"files_id": newImg.FileId}
		err = chunkCol.FindOne(ctx, filter).Decode(&imgChunk)
		newImg.Data = imgChunk.Data
		newimg = &newImg
		tmp = append(tmp, newimg)
			

	}
	m.imgs = tmp

	return m.imgs
}

func (m *ImageManager) Lastest() *Image {
	db := m.client.Database("test")
	collection := db.Collection("test")
	chunkCol := db.Collection("fs.chunks")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var newImg Image
	options := options.Find()
	options.SetSort(bson.D{{"_id", -1}})
	options.SetLimit(1)

	// cur, err := collection.Aggregate
	// count, err := collection.EstimatedDocumentCount(context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// }
	err := collection.FindOne(ctx, bson.D{}).Decode(&newImg)
	// if err != nil {
	// 	log.Fatal(err)
	// }	

	filter := &bson.M{"files_id": newImg.FileId}
	var imgChunk ImgChunk
	err = chunkCol.FindOne(ctx, filter).Decode(&imgChunk)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	if err != nil {
		fmt.Println(err)
	}
	newImg.Data = imgChunk.Data
	var newimg *Image
	newimg = &newImg
	m.imgs = append(m.imgs, newimg)
	return newimg
}


func init(){
	var client *mongo.Client 
	client = InitiateMongoClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successful connected and pinged")
	DefaultImageList = NewImageManager(client)
	
}
