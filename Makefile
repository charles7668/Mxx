export CGO_CFLAGS=-I$(CURDIR)/whisper/include
export CGO_LDFLAGS=-L$(CURDIR)/whisper/lib

build:
	go build .

# use this command with make run ARGS="put your args here"
run:
	go run main.go $(ARGS)

clean:
	go clean
