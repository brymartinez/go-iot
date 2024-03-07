.PHONY: build clean deploy

build:
	go get ./...
	go mod vendor
	env GOOS=linux go build -o -ldflags="-s -w" -o bin/main main.go

clean:
	rm -rf ./bin ./vendor

deploy: clean build
	env SLS_DEBUG=* sls deploy --verbose --stage=local --region=ap-southeast-1