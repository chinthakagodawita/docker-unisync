SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=docker-unisync

VERSION=0.0.1
BUILD=`git rev-parse --short HEAD`

LDFLAGS=-ldflags "-X main.Name=${BINARY} -X main.Version=${VERSION} -X main.Build=${BUILD}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go build ${LDFLAGS} -o ${BINARY}

.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
