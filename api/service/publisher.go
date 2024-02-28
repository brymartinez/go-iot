package service

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

func Publish(clss string, message string) error {
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
		return errors.New("cannot instantiate AWS connection")
	}

	client := sns.NewFromConfig(cfg)

	topicArn := os.Getenv("TOPIC_ARN")

	attributes := map[string]types.MessageAttributeValue{
		"IOT_ACTIVATION": {
			DataType:    aws.String("String"),
			StringValue: aws.String(clss),
		},
	}

	result, err := client.Publish(context.Background(), &sns.PublishInput{
		Message:           &message,
		TopicArn:          &topicArn,
		MessageAttributes: attributes,
	})

	if err != nil {
		fmt.Println("Error publishing message:", err)
		return err
	}

	fmt.Println("Message published to topic:", *result.MessageId)
	return nil
}
