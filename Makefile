vet:
	go vet ./...

test: vet
	go test ./...

release:
	go get github.com/Shyp/bump_version
	bump_version minor main.go
