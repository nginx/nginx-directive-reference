# Reference Library

Machine-readable NGINX directive documentation, with Markdown and HTML formats.

This library is generated from the official NGINX documentation, and intended for use in other tools that work with NGINX configuration files. This provides more accurate and up-to-date information than web scraping <https://nginx.org>.

Every time the NGINX documentation changes, this library will be updated and published with a new version number.

## Installation

`npm install --save @nginx/reference-lib`

## Usage

We export two functions (with type definitions) to navigate the documentation:

1. `getDirectives` returns a flat list of all directives in any official module
2. `find` returns a specific directive by name

See [examples](./examples) for more usage.
