.PHONY:run
run:
	CGO_ENABLED=1 go run main.go

.PHONY:test
test:
	go test ./...

.PHONY:build-plugins
build-plugins:
	CGO_ENABLED=1 go build -buildmode=plugin -o bin/sys.so extensions/plugins/sys/sys.go
	CGO_ENABLED=1 go build -buildmode=plugin -o bin/http.so extensions/plugins/http/http.go