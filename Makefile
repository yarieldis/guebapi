.PHONY: build test build-lambda clean

build:
	go build ./...

test:
	go test ./...

build-lambda:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap ./cmd/lambda
	zip lambda.zip bootstrap
	rm -f bootstrap

clean:
	rm -f lambda.zip bootstrap
