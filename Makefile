# variables
GOCMD			=	go
GOPATH			:=	${shell pwd}
BINPATH			=	$(GOPATH)/bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get


export GOPATH

MAIN_PKGS 		:=	server node

all: deps build

build:
		for target in $(MAIN_PKGS); do \
        	$(GOBUILD) -o $(BINPATH)/$$target ./cmd/$$target; \
        done

deps:
		$(GOGET) github.com/aklyukin/rabbit-test-proto

