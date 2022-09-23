build: main.go parser.go cmd.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./event_bus ./...

run: event_bus
	./event_bus --input-file=./sample.txt

clean: event_bus
	@rm event_bus

debug:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./event_bus -gcflags="all=-N -l" ./...
	@dlv --listen=:2345 --headless=true --api-version=2 exec ./event_bus -- --file=./sample.txt