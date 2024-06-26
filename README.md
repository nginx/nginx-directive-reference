[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)
[![License](https://img.shields.io/badge/License-Apache--2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Community Support](https://badgen.net/badge/support/community/cyan?icon=awesome)](SUPPORT.md)

# NGINX Directive Reference

This repo contains the reference-lib npm package that can be used to fetch NGINX directive reference in Markdown and HTML format. It also has the converter code that is used to generate the directive reference json.

## Getting Started

Refer to the [README file](reference-lib/README.md) for how to install the reference-lib package.

Refer to the [README file](reference-converter/README.md) for how to build and run the reference converter.

This project includes a [devcontainer](https://containers.dev/overview) with all dependencies pre-installed and configured.

## Directory Structure

The repository is organized as follows:

- `reference-lib/`: This directory contains an npm package that allows for easy lookup of nginx directives
- `reference-converter/`: This directory contains a program that converts the official NGINX reference documentation from XML format to JSON.
- `tools/`: This directory contains development tools
- `examples/`: This directory contains example usage of the `reference-lib`

## Contributing

Please see the [contributing guide](CONTRIBUTING.md) for guidelines on how to best contribute to this project.

## License

[Apache License, Version 2.0](LICENSE)

&copy; [F5, Inc.](https://www.f5.com/) 2024
