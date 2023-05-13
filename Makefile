# Image name to build the docker container from
IMAGE=azopat/gomagik

# use the current directory name for the NAME of the container and so on
# relies on always calling the Makefile from the directory it is in but tha is
# the way it was designed to work anyway
NAME=$(notdir $(shell pwd))

# location of the project within the docker VM
ROOT=$(shell pwd)
VENDOR_ROOT=$(ROOT)/../_vendor
COMMON_ROOT=$(ROOT)/../common

IMAGE_BASE_NAME=$(DOCKER_IMAGES_BASE_NAME)/$(NAME)

# location project is mapped to inside the docker container
LOCAL_ROOT=/go/$(NAME)

DOCKER_CMD_TO_RUN=docker run --rm=true

DOCKER_CMD_TO_RUN+= -v $(ROOT):$(LOCAL_ROOT)

DOCKER_CMD_TO_RUN+= -v $(VENDOR_ROOT):/go/_vendor

DOCKER_CMD_TO_RUN+= -v $(COMMON_ROOT):/go/common

DOCKER_CMD_TO_RUN+=  $(IMAGE)

compile: gofmt build
	cp Dockerfile tmp/

build:
	$(DOCKER_CMD_TO_RUN) sh -c "cd /go && rm -rf /go/$(NAME)/bin && GOPATH=/go/_vendor /usr/local/go/bin/go build -ldflags '-X main.BUILD=$(COMMIT_TIME)__$(COMMIT_HASH)'  -o  /go/$(NAME)/tmp/$(NAME) /go/$(NAME)/*.go"

gofmt:
	$(DOCKER_CMD_TO_RUN) sh -c "cd /go && rm -rf /go/$(NAME)/bin && /usr/local/go/bin/gofmt -w /go/$(NAME)/*.go"
