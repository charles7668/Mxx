ifeq ($(OS),Windows_NT)
	export CGO_CFLAGS=-I$(CURDIR)/whisper/include
	export CGO_LDFLAGS=-L$(CURDIR)/whisper/lib/win
    export CC=gcc.exe
    export CXX=g++.exe
    CUR_PATH := $(PATH)
    export PATH=$(CURDIR)/whisper/lib/win:$(CUR_PATH)
else
	export CGO_CFLAGS=-I$(CURDIR)/whisper/include
	export CGO_LDFLAGS=-L$(CURDIR)/whisper/lib
endif

# The mp4 in TestSrc is used from https://www.youtube.com/watch?v=_z6ZIwKu1bY
export FFMPEG_TEST_DIR=$(CURDIR)/TestSrc
export WHISPER_TEST_DIR=$(CURDIR)/TestSrc

.PHONY: build build-backend build-frontend
build-backend:
	go build .
build-frontend:
	cd web && npm install && npm run build
build: build-frontend build-backend

.PHONY: run-api run
# use this command with make run ARGS="put your args here"
run-api:
	go run main.go --api
# this command will run the web server with the frontend
run-web:
	go run main.go --web
# use this command with make run ARGS="put your args here"
run:
	go run main.go $(ARGS)

.PHONY: clean-backend clean-frontend clean
clean-backend:
	go clean
clean-frontend:
	rm -rf $(CURDIR)/dist
clean: clean-backend clean-frontend

.PHONY: test
test:
	go test -v ./...
	rm -rf $(CURDIR)/data/temp
	rm -rf $(CURDIR)/data/media
