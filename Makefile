BINARY_NAME=strava_bot

build:
	GOARCH=amd64 GOOS=darwin go build -o ./output/${BINARY_NAME}-mac cmd/bot/main.go
	GOARCH=amd64 GOOS=linux go build -o ./output/${BINARY_NAME}-linux cmd/bot/main.go
	go build -o ./output/${BINARY_NAME} cmd/bot/main.go

run: build
	./output/${BINARY_NAME}

clean:
	go clean
	rm ./output/${BINARY_NAME}-mac
	rm ./output/${BINARY_NAME}-linux

mock:
	mockery --all

deploy: build
	./deploy.sh