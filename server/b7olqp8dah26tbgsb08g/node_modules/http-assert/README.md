# http-assert

[![NPM Version][npm-image]][npm-url]
[![NPM Downloads][downloads-image]][downloads-url]
[![Node.js Version][node-version-image]][node-version-url]
[![Build Status][travis-image]][travis-url]
[![Test Coverage][coveralls-image]][coveralls-url]

Assert with status codes. Like ctx.throw() in Koa, but with a guard.

## Example
```js
var assert = require('http-assert');
var ok = require('assert');

try {
  assert(username == 'fjodor', 401, 'authentication failed');
} catch (err) {
  ok(err.status == 401);
  ok(err.message == 'authentication failed');
  ok(err.expose);
}
```

## Licence

[MIT](LICENSE)

[npm-image]: https://img.shields.io/npm/v/http-assert.svg
[npm-url]: https://npmjs.org/package/http-assert
[node-version-image]: https://img.shields.io/node/v/http-assert.svg
[node-version-url]: https://nodejs.org/en/download/
[travis-image]: https://img.shields.io/travis/jshttp/http-assert/master.svg
[travis-url]: https://travis-ci.org/jshttp/http-assert
[coveralls-image]: https://img.shields.io/coveralls/jshttp/http-assert/master.svg
[coveralls-url]: https://coveralls.io/r/jshttp/http-assert
[downloads-image]: https://img.shields.io/npm/dm/http-assert.svg
[downloads-url]: https://npmjs.org/package/http-assert
