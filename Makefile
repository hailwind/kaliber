all: build

$(shell if [ ! -d bin ]; then mkdir bin; fi )

x86_64:
	go mod tidy
	@GOOS=linux go build -ldflags '-linkmode "external" -extldflags "-static"' -o bin/kaliber-x86_64 app/kaliber.go

aarch64:
	go mod tidy
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc go build -ldflags '-linkmode "external" -extldflags "-static"' -o bin/kaliber-aarch64 app/kaliber.go

build: x86_64 aarch64

clean:
	rm -f bin/*
