air:
	air

build:
	go build -o bin/main.exe

run: build
	./bin/main

test:
	go test -v ./... -count=1