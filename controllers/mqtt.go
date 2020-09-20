package controllers

import(
	"fmt"
	"os"
	mqtt "github.com/eclipse/paho.mqtt.golang"	
	"myproject/models"
	"encoding/json"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message){
	fmt.Printf("Topic: %s\n", msg.Topic())
	// fmt.Printf("MSG: %s\n", msg.Payload())
	var new_img models.Image
	if err:= json.Unmarshal(msg.Payload(), &new_img); err != nil {
		fmt.Println(err)
	}
	_, err := models.DefaultImageList.InsertImg(new_img)
	if err != nil {
		fmt.Println(err)
	}
}

func Connect() {
	opts := mqtt.NewClientOptions().AddBroker("localhost:1883")
	opts.SetClientID("user1")
	opts.SetDefaultPublishHandler(f)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

	if token:= c.Subscribe("test", 0 , nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}