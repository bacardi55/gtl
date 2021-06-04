GOCMD := CGO_ENABLED=0 go1.16.4
BINARY := gtl
BINDIR := ./bin
VERSION := 0.0.1

GOLDFLAGS ?= -s -w -X git.bacardi55.io/gtl/main.Version=$(VERSION)

#PREFIX := /usr/local
#EXEC_PREFIX := ${PREFIX}
#BINDIR := ${EXEC_PREFIX}/bin
#DATAROOTDIR := ${PREFIX}/share
#MANDIR := ${DATAROOTDIR}/man
#MAN1DIR := ${MANDIR}/man1
#test : GOCMD := go1.11.13

BUILD_TIME := ${shell date "+%Y-%m-%dT%H:%M"}

.PHONY: build
build:
	${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}

.PHONY: clean
clean:
	rm -f ${BINDIR}/${BINARY}

fmt:
	go fmt ./...

.PHONY: release
release:
	echo "Tagging version " ${VERSION}
	git tag v${VERSION}
	git tag push origin v${VERSION}
	GOOS=linux GOARCH=amd64 ${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}_${VERSION}_linux_64
	GOOS=linux GOARCH=arm ${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}_${VERSION}_linux_arm
	GOOS=linux GOARCH=arm64 ${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}_${VERSION}_linux_arm
	GOOS=linux GOARCH=386 ${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}_${VERSION}_linux_32
