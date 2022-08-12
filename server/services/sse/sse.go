package sse

import (
	"encoding/json"
	"github.com/r3labs/sse/v2"
	"log"
)

var SSE = sse.New()

func init() {
	stream := SSE.CreateStream("message")
	stream.AutoReplay = false
}

type Message struct {
	Title string `json:"title"`
	Info  string `json:"info"`
}

func PublishMessage(msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}
	SSE.Publish("message", &sse.Event{
		Data: data,
	})
}
