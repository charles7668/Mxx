export CGO_CFLAGS=-I$(CURDIR)/whisper/include
export CGO_LDFLAGS=-L$(CURDIR)/whisper/lib

# The mp4 in TestSrc is used from https://www.youtube.com/watch?v=_z6ZIwKu1bY
export FFMPEG_TEST_DIR=$(CURDIR)/TestSrc
export WHISPER_TEST_DIR=$(CURDIR)/TestSrc

build-backend:
	go build .
build-frontend:
	cd web && npm install && npm run build
build: build-backend build-frontend

# use this command with make run ARGS="put your args here"
run-api:
	go run main.go $(ARGS)

clean-backend:
	go clean
clean-frontend:
	rm -rf $(CURDIR)/dist
clean: clean-backend clean-frontend

test:
	go test -v ./...
	rm -rf $(CURDIR)/api/data/temp
	rm -rf $(CURDIR)/api/data/media
