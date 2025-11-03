.PHONY: build run
build:
	go build -o ./build/bin/boocol ./cmd/api/
run: build
	./build/bin/boocol
