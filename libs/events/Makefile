.PHONY: docs
REPO:=github.com/vbhp/supermint

docs:
	@go get github.com/davecheney/godoc2md
	godoc2md $(REPO) > README.md

test:
	go test -v ./...
