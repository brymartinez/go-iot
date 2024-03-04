package service

import (
	"context"
	"encoding/json"
	"fmt"
	"go-iot/api/common"
	"go-iot/api/model"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	pg "github.com/go-pg/pg/v10"
)

var db *pg.DB

type SNSEvent struct {
	RequestType       *string                 `json:"Type"`
	SubscribeURL      *string                 `json:"SubscribeURL"`
	Message           *string                 `json:"Message"`
	MessageAttributes *map[string]interface{} `json:"MessageAttributes"`
}

var snsClient *sns.Client

func Subscribe() error {
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
		return err
	}

	snsClient = sns.NewFromConfig(cfg)
	go startHttpServer(":8083")
	subscribeToSNS("http://host.docker.internal:8083")
	return nil
}

func startHttpServer(port string) {
	server := http.Server{
		Addr:    port,
		Handler: handler(),
	}
	log.Println("starting server on port 8083")
	err := server.ListenAndServe()
	log.Fatal(err)
}

func handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}

		fmt.Println("Received body", string(body))
		req := SNSEvent{}

		if err := json.Unmarshal(body, &req); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		log.Println(*req.Message)
		log.Println(*req.RequestType)
		log.Println(*req.MessageAttributes)

		topic := ""
		for key := range *req.MessageAttributes {
			if key == "IOT_ACTIVATION_RESPONSE" || key == "IOT_DEACTIVATION_RESPONSE" {
				topic = key
			}
		}

		//confirmation request
		if *req.RequestType == "SubscriptionConfirmation" {
			confirm(req)
			return
		}

		var message model.Device
		err = json.Unmarshal([]byte(*req.Message), &message)
		if err != nil {
			log.Println("Got string message", *req.Message)
			return
		}

		db, err = common.ConnectToDB()
		if err != nil {
			fmt.Printf("Error connecting to db, %d\n", err)
			w.WriteHeader(500)
			return
		}

		switch topic {
		case "IOT_ACTIVATION_RESPONSE":
			handleActivateResponse(message)
		case "IOT_DEACTIVATION_RESPONSE":
			handleDeactivateResponse(message)
		}

		w.WriteHeader(200)
	}
}

func handleActivateResponse(device model.Device) {
	log.Println("Got object message", device)
	if device.Class == "Other" { // Condition to disapprove "Other" devices
		device.Status = "ACTIVE"
		_, err := db.Model(&device).Where("id = ? AND status='PENDING'", device.ID).Update(&device)
		if err != nil {
			fmt.Printf("Error saving to db, %d\n", err)
			rollback(device)
			return
		}
	} else {
		_, err := db.Model(&device).Where("id = ? AND status='PENDING'", device.ID).Update(&device)
		if err != nil {
			fmt.Printf("Error saving to db, %d\n", err)
			rollback(device)
			return
		}
	}
}

func handleDeactivateResponse(device model.Device) {

}

func publish(message string) {
	log.Println(message)
}

func rollback(device model.Device) {

}

func subscribeToSNS(endpoint string) error {
	topicArn := os.Getenv("TOPIC_ARN")

	filterMap := map[string][]string{
		"IOT_ACTIVATION_RESPONSE":   {"Living Room", "Bedroom", "Dining Room", "Kitchen", "Other"},
		"IOT_DEACTIVATION_RESPONSE": {"Living Room", "Bedroom", "Dining Room", "Kitchen", "Other"},
	}

	var attributes map[string]string
	filterBytes, err := json.Marshal(filterMap)
	if err != nil {
		log.Printf("Couldn't create filter policy, here's why: %v\n", err)
		return err
	}
	attributes = map[string]string{"FilterPolicy": string(filterBytes)}

	protocol := "http"
	// Subscribe to the SNS topic
	output, err := snsClient.Subscribe(context.Background(), &sns.SubscribeInput{
		TopicArn:   &topicArn,
		Protocol:   &protocol,
		Endpoint:   &endpoint,
		Attributes: attributes,
	})

	fmt.Printf("Successful subscription\n%d\n", output.SubscriptionArn)
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

	log.Println("confirm output", *output)
}
