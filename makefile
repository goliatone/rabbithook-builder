BUILDPATH=$(CURDIR)
SRCPATH=$(CURDIR)/src/
GO=$(shell which go)

GOGET=$(GO) get
GOCLEAN=$(GO) clean
GOINSTALL=$(GO) install
GOBUILD=$(GO) build -o
BFLAGS=-ldflags "-w" $(SRCPATH)main.go
EXNAME=rhbuilder

ARCH=mac
BINPATH=$(BUILDPATH)/bin/$(ARCH)/

export GOPATH=$(CURDIR)

myname:
	@echo "I do make files"

makedir:
	@echo "start building tree..."
	@if [ ! -d $(BUILDPATH)/bin ] ; then mkdir -p $(BUILDPATH)/bin ; fi
	@if [ ! -d $(BUILDPATH)/pkg ] ; then mkdir -p $(BUILDPATH)/pkg ; fi

get:
	@$(GOGET) github.com/streadway/amqp

build:
	@echo "Start building..."
	$(GOBUILD) $(BINPATH)$(EXNAME) $(BFLAGS)
	@echo "Yay! all DONE!"

clean:
	@echo "cleanning"
	@rm -rf $(BUILDPATH)/bin/$(EXENAME)
	@rm -rf $(BUILDPATH)/pkg
	@rm -rf $(BUILDPATH)/src/github.com

all: makedir get build
