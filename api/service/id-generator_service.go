package service

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func printMap(m map[string]types.AttributeValue) {
	for key, value := range m {
		fmt.Printf("%s: %v\n", key, value)
	}
}

func IDGenerator(clss string) string {

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
		panic("Cannot instantiate AWS connection.")
	}

	client := dynamodb.NewFromConfig(cfg)

	// Define the parameters for the update operation
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("IDGenerator-local"),
		Key: map[string]types.AttributeValue{
			"IDSequence": &types.AttributeValueMemberS{
				Value: "IDSequence",
			},
		},
		UpdateExpression: aws.String("SET #attrName = #attrName + :value"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":value": &types.AttributeValueMemberN{
				Value: "1",
			},
		},
		ExpressionAttributeNames: map[string]string{
			"#attrName": clss,
		},
		ReturnValues: types.ReturnValueAllNew,
	}

	output, err := client.UpdateItem(context.Background(), input)
	if err != nil {
		log.Fatalf("Got error calling UpdateItem: %s", err)
	} else {
		printMap(output.Attributes)
	}

	var value string
	if attr, found := output.Attributes[clss]; found {
		switch v := attr.(type) {
		case *types.AttributeValueMemberS:
			value = v.Value
		case *types.AttributeValueMemberN:
			value = v.Value
		default:
			log.Fatalf("Unexpected type for attribute %s: %T", clss, attr)
		}
	}

	return value
}
