default: help

# variable collection
# used way: only adjust 【PROJECT】and 【DOCKER_REGISTRY】
PROJECT = main-data-api
HUBPATH = maindata
DOCKER_REGISTRY=docker.osisbim.com
GO_LDFLAGS = -ldflags " -w"
#VERSION = $(shell date -u +v%Y%m%d)-$(shell git describe --tags --always)
VERSION = $(shell date -u +v%Y%m%d)-$(version)
BIN_LABELS = ${PROJECT}_$(VERSION)
WIN_LABELS = ${PROJECT}_$(VERSION).exe
DOCKER_IMAGE_NAME = ${DOCKER_REGISTRY}/${HUBPATH}/${PROJECT}:$(VERSION)
DOCKER_REMOVE_IMAGE_NAME = ${DOCKER_REGISTRY}/${HUBPATH}/${PROJECT}:latest

build-lux:
	echo ${BIN_LABELS}
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${GO_LDFLAGS} -o ${BIN_LABELS} main.go

build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ${GO_LDFLAGS} -o ${BIN_LABELS} main.go

build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${GO_LDFLAGS} -o ${WIN_LABELS} main.go

build-img:
	docker build . -f Dockerfile -t ${DOCKER_IMAGE_NAME} --build-arg BIN_LABELS=${BIN_LABELS}
	docker tag ${DOCKER_IMAGE_NAME} ${DOCKER_REMOVE_IMAGE_NAME}

docker-push: build-img
	docker push ${DOCKER_REMOVE_IMAGE_NAME}
	docker push ${DOCKER_IMAGE_NAME}