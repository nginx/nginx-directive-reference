# autocomplete

This example generates a `dist/directives.json` formatted like the one in used in https://github.com/jaywcjlove/nginx-editor, but instead of web scraping it uses `nginx-directive-reference`.

The result is a little more accurate, with up-to-date information and no false-positives where the web scraper is misinterperting the HTML and documenting non-existent NGINX directives.

## Usage

1. `npm ci`
2. `npm run build`
3. use the file in `dist/directives.json` as needed
4. `npm run highlight`
5. Copy the content from the `dist/highlight.ts` to the `tokenizer: { root: [...]}` section in `monarch-token-provider.ts`
