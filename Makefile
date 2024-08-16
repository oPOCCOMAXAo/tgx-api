.PHONY: build test generate build-docker upload-docker test-docker install-deps

install-deps:
	go install 

generate:
	go generate ./...

build-docker:
	docker build -t poccomaxa/tgx:latest .

upload-docker:
	docker push poccomaxa/tgx:latest

test-docker:
	docker run -it --rm poccomaxa/tgx:latest
