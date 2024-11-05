build:
	go build -ldflags="-s -w" -o bin/baobud .
clean:
	rm -rf bin/
dev:
	go run ./cmd -f test/template.ctmpl
