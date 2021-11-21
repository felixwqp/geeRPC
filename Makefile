hello:
	echo "Hello Go Makefile"

build:
	go build -o bin/main main.go

run:
	go run main/main.go