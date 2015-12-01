SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=docker-unisync

VERSION=1.0.0
BUILD=`git rev-parse --short HEAD`
TARGET_OS="darwin"
TARGET_ARCH="amd64"
BINDIR=bin

LDFLAGS=-ldflags "-X main.Name=${BINARY} -X main.Version=${VERSION} -X main.Build=${BUILD}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go build ${LDFLAGS} -o ${BINDIR}/${BINARY}

.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: bootstrap
bootstrap:
	# Dependencies.
	glide install

.PHONY: clean
clean:
	if [ -f ${BINDIR} ] ; then rm -r ${BINDIR} ; fi
