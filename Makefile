generate-domain-errors:
	ERRORS_YAML_FILE_PATH=api/errors.yaml \
	ERRORS_TARGET_DIR=gen/log \
	ERRORS_TARGET_FILENAME=app-errors.gen.go \
	ERRORS_PACKAGE_NAME=log \
	go run main.go

bin-build:
	go build -o bin/main ./main.go

bin-generate:
	ERRORS_YAML_FILE_PATH=api/errors.yaml \
	ERRORS_TARGET_DIR=gen/log \
	ERRORS_TARGET_FILENAME=app-errors.gen.go \
	ERRORS_PACKAGE_NAME=log \
	bin/main

docker-build:
	docker build -t docker.io/iadolgov/domain-error-go -f Dockerfile.multistage .

docker-deploy:
	docker push iadolgov/domain-error-go

docker-generate:
	docker run --rm -v $(PWD):$(PWD) -w $(PWD) -u `id -u $(USER)` \
	-e ERRORS_YAML_FILE_PATH=api/errors.yaml \
	-e ERRORS_TARGET_DIR=pkg/log2 \
	-e ERRORS_TARGET_FILENAME=app-errors.gen.go \
	-e ERRORS_PACKAGE_NAME=log2 \
	iadolgov/domain-error-go /app/run

lint:
	docker run --rm -v $(PWD):$(PWD) -w $(PWD) -u `id -u $(USER)` \
	-e GOLANGCI_LINT_CACHE=/tmp/.cache \
	-e GOCACHE=/tmp/.cache golangci/golangci-lint:v1.61.0 \
	golangci-lint run -v --fix

test:
	go test ./...
