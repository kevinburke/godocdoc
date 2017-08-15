BUMP_VERSION := $(GOPATH)/bin/bump_version
GODOCDOC := $(GOPATH)/bin/godocdoc
GO_FILES := $(shell find . -name '*.go')
MEGACHECK := $(GOPATH)/bin/megacheck

$(BUMP_VERSION):
	go get -u github.com/Shyp/bump_version

$(MEGACHECK):
	go get -u honnef.co/go/tools/cmd/megacheck

lint: $(MEGACHECK)
	$(MEGACHECK) ./...
	go vet ./...

test: lint
	go test ./...

release: $(BUMP_VERSION)
	$(BUMP_VERSION) minor main.go

$(GODOCDOC): $(GO_FILES)
	go install -v .

serve: $(GODOCDOC)
	$(GODOCDOC)
