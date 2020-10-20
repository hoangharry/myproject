package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"myproject/models"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	// "strings"
)

type RawImg struct {
	Usr       string `json:"usr" bson:"usr"`
	CamID     string `json:"cam_id" bson:"cam_id"`
	Line      int    `json:"line" bson:"line"`
	Direction int    `json:"direction" bson:"direction"`
	Time      string `json:"time" bson:"time"`
	Evidence  string `json:"evidence" bson:"evidence"`
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var rawimg RawImg
	if err := json.Unmarshal(msg.Payload(), &rawimg); err != nil {
		fmt.Println(err)
	}
	fmt.Println(msg.Topic())
	tmp, _ := time.Parse(time.RFC3339, rawimg.Time)
	newImg := models.Image{Usr: rawimg.Usr, CamID: rawimg.CamID, Line: rawimg.Line, Direction: rawimg.Direction, Time: tmp, Evidence: rawimg.Evidence}

	go AlterImg(newImg)

	go BroadcastWs(newImg)
}

func AlterImg(newImg models.Image) {
	unbase, err := base64.StdEncoding.DecodeString(newImg.Evidence[2 : len(newImg.Evidence)-1])
	if err != nil {
		fmt.Println(err)
	}
	r := bytes.NewReader(unbase)
	img, err := jpeg.Decode(r)
	if err != nil {
		fmt.Println(err)
	}
	filename := "LocalStorage/" + newImg.Usr + newImg.Time.String() + ".jpeg"
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println(err)
	}
	jpeg.Encode(f, img, nil)
	newImg.Evidence = filename
	fmt.Println(newImg)
	_, err = models.DefaultImageList.InsertImg(newImg)
	if err != nil {
		fmt.Println(err)
	}
}

func Connect() {
	opts := mqtt.NewClientOptions().AddBroker("broker.emqx.io:1883")
	// opts := mqtt.NewClientOptions().AddBroker("localhost:1883")
	opts.SetClientID("user1")
	opts.SetDefaultPublishHandler(f)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

	if token := c.Subscribe("hmh-pga-internship-2020", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

}
