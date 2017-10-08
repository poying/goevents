test: mocks
	@go test .

requirements:
	@go get github.com/golang/mock/mockgen

mocks:
	@-mkdir -p "$@"
	@mockgen -destination "$@/mocks.go" github.com/poying/goevents Producer,Consumer,EventHandler,Message

clean:
	@rm -rf mocks

.PHONY: mocks