package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Report struct {
	Usr       string    `json:"usr" bson:"usr"`
	Time      time.Time `json:"time" bson:"time"`
	CamID     string    `json:"cam_id" bson:"cam_id"`
	Line      int       `json: "line" bson:"line"`
	Direction int       `json:"direction" bson:"direction"`
	Total     int       `json:"total" bson:"total"`
}

type ReportManager struct {
	reports []*Report
	client  *mongo.Client
}

func NewReportManager(client *mongo.Client) *ReportManager {
	return &ReportManager{client: client}
}

var DefaultReportManager *ReportManager

func (m *ReportManager) FindReportOnUsr(fromDate time.Time, toDate time.Time, usr string) []*Report {
	collection := m.client.Database("CloudCam").Collection("Reports")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	option := options.Find()
	option.SetSort(bson.D{{"time", -1}})
	cursor, err := collection.Find(ctx, bson.M{
		"time": bson.M{"$gte": fromDate, "$lte": toDate}, "usr": usr}, option)

	if err != nil {
		fmt.Println(err)
	}
	defer cursor.Close(ctx)
	var tmp []*Report
	for cursor.Next(ctx) {
		var newrp *Report
		var newRp Report
		err = cursor.Decode(&newRp)
		if err != nil {
			fmt.Println(err)
		}
		newrp = &newRp
		tmp = append(tmp, newrp)
	}
	m.reports = tmp

	return m.reports
}

func (m *ReportManager) FindReportOnCam(fromDate time.Time, toDate time.Time, usr string, camID string) []*Report {
	collection := m.client.Database("CloudCam").Collection("Reports")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{
		"time": bson.M{"$gte": fromDate, "$lte": toDate}, "usr": usr, "cam_id": camID})

	if err != nil {
		fmt.Println(err)
	}
	defer cursor.Close(ctx)
	var tmp []*Report
	for cursor.Next(ctx) {
		var newrp *Report
		var newRp Report
		err = cursor.Decode(&newRp)
		if err != nil {
			fmt.Println(err)
		}
		newrp = &newRp
		tmp = append(tmp, newrp)
	}
	if toDate.Format("01-01-2006") == time.Now().Format("01-01-2006") {
		imgCol := m.client.Database("CloudCam").Collection("Data")
		cursor, err = imgCol.Find(ctx, bson.M{"time": bson.M{"$gte": toDate}, "usr": usr, "cam_id": camID})
		if err != nil {
			fmt.Println(err)
		}
		defer cursor.Close(ctx)
		var todayReport0 = Report{Usr: usr, CamID: camID, Time: toDate, Line: 0, Direction: 0, Total: 0}
		var todayReport1 = Report{Usr: usr, CamID: camID, Time: toDate, Line: 0, Direction: 1, Total: 0}
		var todayReport2 = Report{Usr: usr, CamID: camID, Time: toDate, Line: 1, Direction: 0, Total: 0}
		var todayReport3 = Report{Usr: usr, CamID: camID, Time: toDate, Line: 1, Direction: 1, Total: 0}
		for cursor.Next(ctx) {
			var img Image
			err = cursor.Decode(&img)
			if err != nil {
				fmt.Println(err)
			}
			if img.Line == 0 {
				if img.Direction == 0 {
					todayReport0.Total++
				} else {
					todayReport1.Total++
				}
			} else {
				if img.Direction == 0 {
					todayReport2.Total++
				} else {
					todayReport3.Total++
				}
			}
		}
		tmp = append(tmp, &todayReport0)
		tmp = append(tmp, &todayReport1)
		tmp = append(tmp, &todayReport2)
		tmp = append(tmp, &todayReport3)
	}
	m.reports = tmp

	return m.reports
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
	DefaultReportManager = NewReportManager(client)
}
