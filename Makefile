GOCMD := CGO_ENABLED=0 go
BINARY := gtl
BINDIR := ./bin
VERSION := 0.5.0

GOLDFLAGS := -s -w -X main.Version=$(VERSION)

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
	echo "Tagging version ${VERSION}"
	git tag -a v${VERSION} -m "New released tag: v${VERSION}"
	GOOS=linux GOARCH=amd64 ${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm ${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}_${VERSION}_linux_arm
	GOOS=linux GOARCH=arm64 ${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}_${VERSION}_linux_arm64
	GOOS=linux GOARCH=386 ${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}_${VERSION}_linux_386

.PHONY: dependencies
dependencies:
	${GOCMD} get "git.sr.ht/~adnano/go-gemini"
	${GOCMD} get "github.com/spf13/pflag"
	${GOCMD} get "github.com/fatih/color"
	${GOCMD} get "github.com/mitchellh/go-homedir"
	${GOCMD} get "github.com/pelletier/go-toml"
	${GOCMD} get "code.rocketnine.space/tslocum/cview"
	${GOCMD} get "github.com/gdamore/tcell"
	${GOCMD} get "code.rocketnine.space/tslocum/cbind"
