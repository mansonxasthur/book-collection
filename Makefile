.PHONY: build build-cli run migrate
build:
	go build -o ./build/bin/boocol ./cmd/api/
build-cli:
	go build -o ./build/bin/boocol-cli ./cmd/cli/
run: build
	./build/bin/boocol
migrate: build-cli
	./build/bin/boocol-cli migration migrate
