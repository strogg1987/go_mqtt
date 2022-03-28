package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mattn/go-tty"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("go_test_client")
	opts.SetKeepAlive(60 * time.Second)
	// Set the message callback handler
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	// Subscribe to a topic
	if token := c.Subscribe("testtopic/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for {
		r, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		//send message on key "z"
		if string(r) == "z" {
			token := c.Publish("testtopic/1", 0, false, "Hello World")
			token.Wait()
			if token.Error() != nil {
				fmt.Println(token.Error())
			}
		}
	}
}
