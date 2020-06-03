.PHONY: build clean deploy

test:
	go test ../...

build: test
	env GOOS=linux go build -ldflags="-s -w" -o bin/admin function/admin/*.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/survey function/survey/*.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
