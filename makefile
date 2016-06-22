BUILDPATH=$(CURDIR)
SRCPATH=$(CURDIR)/src/
GO=$(shell which go)

GOGET=$(GO) get
GOCLEAN=$(GO) clean
GOINSTALL=$(GO) install
GOBUILD=$(GO) build -o

#TODO: include version information in the build
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

goget:
	@$(GOGET) github.com/streadway/amqp

build:
	@echo "Start building..."
	$(GOBUILD) $(BINPATH)$(EXNAME) $(BFLAGS)
	@echo "Yay! all DONE!"

remote:
	# ssh root@wee-1.local "cd /root/CODE/rabbithook-builder/src; make ARCH=arm all"
	@scp root@wee-1.local:/root/CODE/rabbithook-builder/bin/arm/rhbuilder $(BUILDPATH)/bin/arm/

clean:
	@echo "cleanning"
	@rm -rf $(BUILDPATH)/bin/$(EXENAME)
	@rm -rf $(BUILDPATH)/src/github.com

all: makedir goget build
