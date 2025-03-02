build:
	go build -o binary/p2p

run: build
	./binary/p2p -name sai -port 8080 -local true

test:
	go test ./... -v