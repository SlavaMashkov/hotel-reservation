air:
	air

seed:
	go run scripts/seed.go

build:
	go build -o bin/main.exe

run: build
	./bin/main

mockery:
	mockery

test: mockery
	go test -v ./... -count=1