BINARY_NAME=rings-backend
IMAGE_NAME=ghcr.io/denysvitali/rings-social-backend
IMAGE_TAG=latest

build:
	CGO_ENABLED=0 go build -o build/$(BINARY_NAME) ./


docker-build:
	docker build \
		-t "$(IMAGE_NAME):$(IMAGE_TAG)" \
		.

docker-push:
	docker push "$(IMAGE_NAME):$(IMAGE_TAG)"