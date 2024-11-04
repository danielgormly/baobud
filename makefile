build:
	go build -ldflags="-s -w" -o bin/app main.go
clean:
	rm -rf bin/
