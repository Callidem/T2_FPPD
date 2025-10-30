
.PHONY: all build server client clean distclean

all: build

go.mod:
	go mod init jogo
	go get -u github.com/nsf/termbox-go

build: go.mod server client

server:
	go build -o bin/server ./cmd/server

client:
	go build -o bin/cliente ./cmd/client

clean:
	rm -f bin/server bin/cliente jogo

distclean: clean
	rm -f go.mod go.sum
