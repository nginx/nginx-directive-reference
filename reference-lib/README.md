# Reference Library
NGINX directive reference in Markdown and HTML format

# Installation

1. Generate a github personal access token with read:packages permission
2. Create a ~/.npmrc file with the following content
```
@nginxinc:registry=http://npm.pkg.github.com
//npm.pkg.github.com/:_authToken=$TOKEN
```
3. Run
```bash
npm install --save @nginx/reference-lib@1.0.0
```

# Usage
1. find
```javascript
import { find, Format } from '@nginx/reference-lib'
const content = find('listen', undefined, Format.HTML)
```

2. getDirectives
```javascript
import { getDirectives, Format } from '@nginx/reference-lib'
const directive = getDirectives(Format.HTML)
```
