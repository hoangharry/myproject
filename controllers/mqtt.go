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
	fmt.Println(new_img.ID)
	fmt.Println(new_img.Time)
	_, err := models.DefaultImageList.InsertImg(new_img)
	if err != nil {
		fmt.Println(err)
	}
	BroadcastWs(new_img)
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

	if token:= c.Subscribe("hmh-pga-internship-2020", 0 , nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}


}