export CGO_CFLAGS=-I$(CURDIR)/whisper/include
export CGO_LDFLAGS=-L$(CURDIR)/whisper/lib

# The mp4 in TestSrc is used from https://www.youtube.com/watch?v=_z6ZIwKu1bY
export FFMPEG_TEST_DIR=$(CURDIR)/TestSrc
export WHISPER_TEST_DIR=$(CURDIR)/TestSrc

build:
	go build .

# use this command with make run ARGS="put your args here"
run:
	go run main.go $(ARGS)

clean:
	go clean

test:
	go test ./... -v
