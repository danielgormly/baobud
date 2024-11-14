.PHONY: test

build:
	go build -ldflags="-s -w" -o bin/baobud cmd/main.go cmd/file.go
clean:
	rm -rf bin/
dev:
	go run ./cmd -f test/template.ctmpl
start-bao:
	bao server -dev -dev-root-token-id=dev > ./bao.log 2>&1 & echo $$! > vault.pid
	sleep 2
stop-bao:
	@if [ -f vault.pid ]; then \
		kill `cat vault.pid` > /dev/null 2>&1 || true; \
		rm vault.pid; \
	fi
test-binary:
	./bin/baobud version && \
	./bin/baobud -f test/template.ctmpl && \
	kill $$VAULT_PID
test:
	go test ./core
