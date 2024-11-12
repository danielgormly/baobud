build:
	go build -ldflags="-s -w" -o bin/baobud cmd/main.go cmd/file.go
clean:
	rm -rf bin/
dev:
	go run ./cmd -f test/template.ctmpl -o policy.hcl
test-binary:
	bao server -dev & \
	VAULT_PID=$$! && \
	sleep 2 && \
	./bin/baobud version && \
	./bin/baobud -f test/template.ctmpl && \
	kill $$VAULT_PID
