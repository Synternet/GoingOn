.PHONY: build

build:
	go build -o . ./...
build-docker:
	docker build -f ./docker/Dockerfile -t goingon .
run-docker:
	docker run -it --rm --env-file=.env goingon
format:
	gofumpt -l -w .
