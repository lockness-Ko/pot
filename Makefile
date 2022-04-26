all:
	docker build -t minimal .
run: all
	go run .