package models

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Image struct {
	Usr       string    `json:"usr" bson:"usr"`
	CamID     string    `json:"cam_id" bson:"cam_id"`
	Direction int       `json:"direction" bson:"direction"`
	Time      time.Time `json:"time" bson:"time"`
	Line      int       `json:"line" bson:"line"`
	Evidence  string    `json:"evidence" bson:"evidence"`
}

type ImageManager struct {
	imgs   []*Image
	client *mongo.Client
}

var DefaultImageList *ImageManager

func NewImageManager(client *mongo.Client) *ImageManager {
	return &ImageManager{client: client}
}

func (m *ImageManager) PageImg(usr string, camID string, page int) []*Image {
	collection := m.client.Database("CloudCam").Collection("Data")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	option := options.Find()
	option.SetLimit(10)
	option.SetSkip(int64((page - 1) * 10))
	option.SetSort(bson.D{{"time", -1}})
	cursor, err := collection.Find(ctx, bson.M{"usr": usr, "cam_id": camID}, option)
	if err != nil {
		fmt.Println(err)
	}
	defer cursor.Close(ctx)
	var tmp []*Image
	for cursor.Next(ctx) {
		var newImg Image
		err = cursor.Decode(&newImg)
		if err != nil {
			fmt.Println(err)
		}
		f, _ := os.Open(newImg.Evidence)
		reader := bufio.NewReader(f)
		content, _ := ioutil.ReadAll(reader)
		newImg.Evidence = base64.StdEncoding.EncodeToString(content)

		tmp = append(tmp, &newImg)

	}
	m.imgs = tmp
	return m.imgs
}

func (m *ImageManager) CountImg(usr string, camID string) int64 {
	collection := m.client.Database("CloudCam").Collection("Data")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	count, err := collection.CountDocuments(ctx, bson.M{"usr": usr, "cam_id": camID})
	if err != nil {
		fmt.Println(err)
	}
	return count
}

func (m *ImageManager) CountImgMonitor(usr string, camID string, fromTime time.Time) int64 {
	collection := m.client.Database("CloudCam").Collection("Data")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	count, err := collection.CountDocuments(ctx, bson.M{"usr": usr, "cam_id": camID, "time": bson.M{"$gte": fromTime}})
	if err != nil {
		fmt.Println(err)
	}
	return count
}

func (m *ImageManager) CountImgFromQuery(usr string, camID string, fromTime time.Time, toTime time.Time) int64 {
	collection := m.client.Database("CloudCam").Collection("Data")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	count, err := collection.CountDocuments(ctx, bson.M{"usr": usr, "cam_id": camID, "time": bson.M{"$gte": fromTime, "$lte": toTime}})
	if err != nil {
		fmt.Println(err)
	}
	return count
}

func (m *ImageManager) ImgOnInit(usr string, camID string, fromTime time.Time) []*Image {
	collection := m.client.Database("CloudCam").Collection("Data")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	option := options.Find()
	option.SetSort(bson.D{{"time", -1}})
	cursor, err := collection.Find(ctx, bson.M{"usr": usr, "cam_id": camID, "time": bson.M{"$gte": fromTime}}, option)
	if err != nil {
		fmt.Println(err)
	}
	defer cursor.Close(ctx)
	var tmp []*Image
	for cursor.Next(ctx) {
		var newImg Image
		err = cursor.Decode(&newImg)
		if err != nil {
			fmt.Println(err)
		}
		f, _ := os.Open(newImg.Evidence)
		reader := bufio.NewReader(f)
		content, _ := ioutil.ReadAll(reader)
		newImg.Evidence = base64.StdEncoding.EncodeToString(content)

		tmp = append(tmp, &newImg)

	}
	m.imgs = tmp

	return m.imgs
}

func (m *ImageManager) ImgInTime(usr string, camID string, fromTime time.Time, toTime time.Time, page int) []*Image {
	collection := m.client.Database("CloudCam").Collection("Data")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	option := options.Find()
	option.SetLimit(10)
	option.SetSkip(int64((page - 1) * 10))
	option.SetSort(bson.D{{"time", -1}})
	cursor, err := collection.Find(ctx, bson.M{"usr": usr, "cam_id": camID, "time": bson.M{"$gte": fromTime, "$lte": toTime}}, option)
	if err != nil {
		fmt.Println(err)
	}
	defer cursor.Close(ctx)
	var tmp []*Image
	for cursor.Next(ctx) {
		var newImg Image
		err = cursor.Decode(&newImg)
		if err != nil {
			fmt.Println(err)
		}
		f, _ := os.Open(newImg.Evidence)
		reader := bufio.NewReader(f)
		content, _ := ioutil.ReadAll(reader)
		newImg.Evidence = base64.StdEncoding.EncodeToString(content)

		tmp = append(tmp, &newImg)

	}
	m.imgs = tmp

	return m.imgs
}

func (m *ImageManager) InsertImg(img Image) (bool, error) {
	collection := m.client.Database("CloudCam").Collection("Data")
	_, err := collection.InsertOne(context.TODO(), img)
	if err != nil {
		return false, err
	}
	return true, nil
}

func init() {
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
