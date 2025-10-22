
.PHONY: all build server client clean distclean

all: build

go.mod:
	go mod init jogo
	go get -u github.com/nsf/termbox-go

build: go.mod server client

server:
	go build -o server ./cmd/server

client:
	go build -o cliente ./cmd/client

clean:
	rm -f jogo server cliente

distclean: clean
	rm -f go.mod go.sum
