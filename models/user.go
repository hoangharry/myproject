package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type User struct {
	Usrname string   `json:"usrname" bson:"usrname"`
	Pwd     string   `json:"pwd" bson:"pwd"`
	Camera  []string `json:"cam" bson:"cam"`
}

type UserManager struct {
	user   *User
	client *mongo.Client
}

var DefaultUserManager *UserManager

func NewUserManager(client *mongo.Client) *UserManager {
	return &UserManager{client: client}
}

func init() {
	client := InitiateMongoClient()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successful connected and pinged")
	DefaultUserManager = NewUserManager(client)
}

func (m *UserManager) CheckUser(usrname string, pwd string) bool {
	var user *User
	user = m.GetUser(usrname)
	if pwd != user.Pwd {
		return false
	}
	m.user = user
	return true
}

// func (m *UserManager) AddUser(usr string, pwd string) bool {
// 	var user *User
// 	collection := m.client.Database("test").Collection("users")
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	user = &User{Usrname: usr,Pwd: pwd}
// 	m.user = append(m.user, user)
// 	insertRes, err := collection.InsertOne(ctx, user)
// 	if (err != nil) {
// 		log.Fatal(err)
// 		return false
// 	}
// 	if insertRes != nil {
// 		return true
// 	}

// 	return true
// }

func (m *UserManager) GetUser(usrname string) *User {
	db := m.client.Database("test")
	collection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user User
	err := collection.FindOne(ctx, bson.M{"usrname": usrname}).Decode(&user)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &user
}
func (m *UserManager) GetCam(usr string) []string {
	collection := m.client.Database("CloudCam").Collection("Users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user User
	err := collection.FindOne(ctx, bson.M{"usr": usr}).Decode(&user)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return user.Camera
}
