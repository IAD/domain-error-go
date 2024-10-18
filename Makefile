generate:
	ERRORS_YAML_FILE_PATH=api/errors.yaml \
	ERRORS_TARGET_DIR=gen/log \
	ERRORS_TARGET_FILENAME=app-errors.gen.go \
	ERRORS_PACKAGE_NAME=log \
	go run main.go
	goimports -w gen/log/app-errors.gen.go

lint:
	docker run --rm -v $(PWD):$(PWD) -w $(PWD) -u `id -u $(USER)` \
	-e GOLANGCI_LINT_CACHE=/tmp/.cache \
	-e GOCACHE=/tmp/.cache golangci/golangci-lint:v1.61.0 \
	golangci-lint run -v --fix

test:
	go test -count=1 ./...

public-docker-generate:
	docker pull ghcr.io/iad/domain-error-go:latest
	docker run --rm -v $(PWD):$(PWD) -w $(PWD) -u `id -u $(USER)` \
	-e ERRORS_YAML_FILE_PATH=api/errors.yaml \
	-e ERRORS_TARGET_DIR=gen/log2 \
	-e ERRORS_TARGET_FILENAME=app-errors.gen.go \
	-e ERRORS_PACKAGE_NAME=log2 \
	ghcr.io/iad/domain-error-go /app/run
	goimports -w gen/log2/app-errors.gen.go
