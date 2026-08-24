.PHONY: build dev test clean

build:
	npm run build
	go build -trimpath -ldflags="-s -w" -o bin/jmainguy.com .

dev:
	npm run build
	go run .

test:
	npm run build
	go test ./...

clean:
	rm -rf bin web/dist
