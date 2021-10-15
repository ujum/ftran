.PHONY: build run compile clean
BUILD_DIR=./out
BINARY_NAME=${BUILD_DIR}/ftran
SOURCE_MAIN_NAME=./cmd/ftrancli/main.go

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${BINARY_NAME} ${SOURCE_MAIN_NAME}

compile:
	# 64-Bit
	# Linux
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${BINARY_NAME}-linux-amd64.bin ${SOURCE_MAIN_NAME}
	# Windows
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ${BINARY_NAME}-windows-amd64.exe ${SOURCE_MAIN_NAME}

test:
	go test ./... -cover

clean-build:
	#go clean
	rm -rfv ${BUILD_DIR}

test-cover-report:
	mkdir -p ${BUILD_DIR}
	go test ./... -cover -coverprofile=${BUILD_DIR}/test-coverage.out
	go tool cover -html=${BUILD_DIR}/test-coverage.out -o ${BUILD_DIR}/test-coverage.html
