package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNSEvent struct {
	RequestType  *string `json:"Type"`
	SubscribeURL *string `json:"SubscribeURL"`
	Message      *string `json:"Message"`
}

var snsClient *sns.Client

func startHttpServer(wg *sync.WaitGroup, port string) {
	defer wg.Done()
	server := http.Server{
		Addr:    port,
		Handler: handler(),
	}
	log.Println("starting server on port 8080")
	err := server.ListenAndServe()
	log.Fatal(err)
}

func handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("Successful handler event\n%d\n", body)

		req := SNSEvent{}

		if err := json.Unmarshal(body, &req); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		//confirmation request
		if *req.RequestType == "SubscriptionConfirmation" {
			confirm(req)
			return
		}

		//messages
		log.Println("Got Message", *req.Message)
		w.WriteHeader(200)
	}
}

func subscribeToSNS(httpURL string, client *sns.Client) error {
	topicArn := "arn:aws:sns:ap-southeast-1:000000000000:GO_IOT"

	endpoint := "http://host.docker.internal:8082"
	protocol := "http"
	// Subscribe to the SNS topic
	output, err := client.Subscribe(context.Background(), &sns.SubscribeInput{
		TopicArn: &topicArn,
		Protocol: &protocol,
		Endpoint: &endpoint,
	})

	fmt.Printf("Successful subscription\n%d\n", output)
	if err != nil {
		return err
	}

	return nil

}

func confirm(request SNSEvent) {

	vals, err := url.ParseQuery(*request.SubscribeURL)
	if err != nil {
		log.Println(err)
	}

	token := vals.Get("Token")
	topicARN := vals.Get("TopicArn")

	confirm := &sns.ConfirmSubscriptionInput{
		Token:    aws.String(token),
		TopicArn: aws.String(topicARN),
	}
	output, err := snsClient.ConfirmSubscription(context.Background(), confirm)
	if err != nil {
		log.Println(err)
	}

	log.Println("confirm output", output)
}

func main() {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("ap-southeast-1"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:4566"}, nil
			}),
		),
	)

	if err != nil {
		log.Fatal("Error loading AWS config:", err)
	}

	snsClient = sns.NewFromConfig(cfg)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go startHttpServer(wg, ":8082")
	subscribeToSNS("http://host.docker.internal:8082", snsClient)

	wg.Wait()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt
	fmt.Println("Shutting down...")
}
