# Meta info
NAME:=wasmdemo
VERSION:=0.0.1
REVISION:=$(shell git rev-parse --short HEAD)
LDFLAGS:= -X 'main.version=$(VERSION)' \
	-X 'main.revision=$(REVISION)'

## 必要なツール群をセットアップする
setup:
	go get github.com/stretchr/testify/assert
	go get github.com/fogleman/gg

## run test
test:
	go test -v -cover ./...

## run test
test-watch:
	gomon -R -t

## run test
build-watch:
	gomon -R -- make build


## Linitng
lint: setup
	go vet ./...

## formatting 
fmt: setup
	goimports -w ./...

## build 
build:
	GOOS=js GOARCH=wasm go build -o desc/test.wasm main.go

## bench 
bench:
	go test -benchmem ./examples -bench ^Bench

.PHONY: setup deps update test lint


