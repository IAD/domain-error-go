# Domain Error Generator

This project provides a tool for generating domain error definitions in Go from a YAML file. It includes commands for building, generating, linting, and testing the application.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Makefile Commands](#makefile-commands)
- [Contributing](#contributing)
- [License](#license)

## Usage

You can use the following commands defined in the Makefile
```shell
	docker run --rm -v $(PWD):$(PWD) -w $(PWD) -u `id -u $(USER)` \
	-e ERRORS_YAML_FILE_PATH=api/errors.yaml \
	-e ERRORS_TARGET_DIR=pkg/log2 \
	-e ERRORS_TARGET_FILENAME=app-errors.gen.go \
	-e ERRORS_PACKAGE_NAME=log2 \
	iadolgov/domain-error-go /app/run
```

## Contributing

Contributions are welcome! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -m 'Add some feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Open a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.