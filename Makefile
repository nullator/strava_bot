BINARY_NAME=strava_bot

build:
	GOARCH=amd64 GOOS=darwin go build -o ./output/${BINARY_NAME}-mac cmd/bot/main.go
	GOARCH=amd64 GOOS=linux go build -o ./output/${BINARY_NAME}-linux cmd/bot/main.go

run: build
	./output/${BINARY_NAME}-mac

clean:
	go clean
	rm ./output/${BINARY_NAME}-mac
	rm ./output/${BINARY_NAME}-linux
