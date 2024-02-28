#!/usr/bin/env bash
awslocal dynamodb create-table \
    --table-name IDGenerator-local \
    --key-schema AttributeName=IDSequence,KeyType=HASH \
    --attribute-definitions AttributeName=IDSequence,AttributeType=S \
    --billing-mode PAY_PER_REQUEST \
    --region ap-southeast-1

awslocal dynamodb put-item --table-name IDGenerator-local --item '{"IDSequence": { "S": "IDSequence" },"Living Room": { "N": "10000000" },"Bedroom": { "N": "20000000" },"Dining Room": { "N": "30000000" },"Kitchen": { "N": "40000000" },"Other": { "N": "50000000" }}' --region ap-southeast-1
awslocal sns create-topic --name GO_IOT --region ap-southeast-1