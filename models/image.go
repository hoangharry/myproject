package models

import(
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"context"
	"time"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)


type Image struct {
	ID  int `json:"_id" bson:"_id"`
	Direction string `json:"direction" bson:"direction"`
	Time string `json:"time" bson:"time"`
	Line int `json:"line" bson:"line"`
	Evidence string `json:"evidence" bson:"evidence"`
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
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println(err)
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
		newImg.Evidence = newImg.Evidence[2: len(newImg.Evidence) - 1]
		newimg = &newImg
		tmp = append(tmp, newimg)
			

	}
	m.imgs = tmp

	return m.imgs
}

func (m *ImageManager) Lastest() *Image {
	db := m.client.Database("test")
	collection := db.Collection("test")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var newImg Image
	option := options.FindOne()
	option.SetSort(bson.D{{"_id", -1}})

	err := collection.FindOne(ctx, bson.D{}, option).Decode(&newImg)
	if err != nil {
		fmt.Println(err)
	}	
	var newimg *Image
	newimg = &newImg
	m.imgs = append(m.imgs, newimg)
	return newimg
}

func (m *ImageManager) PageOldest() []*Image {
	db := m.client.Database("test")
	collection := db.Collection("test")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	option := options.Find()
	option.SetLimit(30)
	cursor, err := collection.Find(ctx, bson.D{}, option)
	if err != nil {
		fmt.Println(err)
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

		newimg = &newImg
		tmp = append(tmp, newimg)
			

	}
	m.imgs = tmp

	return m.imgs
}

func (m *ImageManager) PageNewest() []*Image {
	db := m.client.Database("test")
	collection := db.Collection("test")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	option := options.Find()
	option.SetSort(bson.D{{"_id", -1}})
	option.SetLimit(30)
	cursor, err := collection.Find(ctx, bson.D{}, option)
	if err != nil {
		fmt.Println(err)
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

		newimg = &newImg
		tmp = append(tmp, newimg)
			

	}
	m.imgs = tmp

	return m.imgs
}

// func (m *ImageManager) InsertImg(id int, img_time string, direction string, line int, evidence string) (bool, error) {
// 	db := m.client.Database("test")
// 	collection := db.Collection("test")
// 	new_img := Image{ID: id, Time: img_time, Direction: direction, Line: line, Evidence: evidence}
// 	insertRes, err := collection.InsertOne(context.TODO(), new_img)
// 	if err != nil {
// 		fmt.Print(err)
// 		return false, err
// 	}
// 	fmt.Println(insertRes.InsertedID)
// 	return true, nil
	
// }
func (m *ImageManager) InsertImg(img Image) (bool, error){
	collection := m.client.Database("test").Collection("test")
	_, err := collection.InsertOne(context.TODO(), img)
	if err != nil {
		return false, err
	}
	return true, nil
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
