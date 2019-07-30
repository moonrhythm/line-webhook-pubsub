package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/line/line-bot-sdk-go/linebot"
	"gocloud.dev/pubsub"
)

var (
	port              string
	lineChannelSecret string
	pubsubURL         string
)

var (
	topic *pubsub.Topic
	bgCtx = context.Background()
)

func loadEnv() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	lineChannelSecret = os.Getenv("LINE_CHANNEL_SECRET")
	if lineChannelSecret == "" {
		log.Fatal("no LINE_CHANNEL_SECRET provided")
	}

	pubsubURL = os.Getenv("PUBSUB_URL")
	if pubsubURL == "" {
		log.Fatal("no PUBSUB_URL provided")
	}
}

func main() {
	loadEnv()

	var err error
	topic, err = pubsub.OpenTopic(bgCtx, pubsubURL)
	if err != nil {
		log.Fatal("open topic error; ", err)
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go log.Fatal(http.ListenAndServe(":"+port, http.HandlerFunc(handler)))

	<-stop

	err = topic.Shutdown(bgCtx)
	if err != nil {
		log.Fatal("shutdown topic error; ", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	events, err := linebot.ParseRequest(lineChannelSecret, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, event := range events {
		body, err := json.Marshal(event)
		if err != nil {
			log.Println("marshal event error;", err)
			continue
		}
		err = topic.Send(bgCtx, &pubsub.Message{
			Body: body,
		})
		if err != nil {
			log.Println("sending event error;", err)
		}
	}
}
