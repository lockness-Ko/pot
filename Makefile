docker:
	docker build -t minimal .
run:
	go run .
gobuild: clean
	export CGO_ENABLED=0
	export GOOS=linux
	go build -ldflags="-extldflags=-static" -buildmode=pie .
	strip ./pot
clean:
	-rm ./pot
