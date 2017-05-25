BUMP_VERSION := $(shell command -v bump_version)

vet:
	go vet ./...

test: vet
	go test ./...

release:
ifndef BUMP_VERSION
	go get -u github.com/Shyp/bump_version
endif
	bump_version minor main.go
