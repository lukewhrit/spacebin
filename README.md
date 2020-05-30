# Glue

[![Requires.io](https://img.shields.io/requires/github/324Luke/glue)](https://requires.io/github/324Luke/glue/requirements/?branch=master) [![GitHub last commit](https://img.shields.io/github/last-commit/324Luke/glue)](https://github.com/324Luke/glue/commits/master) [![TravisCI](https://img.shields.io/travis/324Luke/glue)](https://travis-ci.org/github/324Luke/glue)

>  Glue is a modern pastebin service, built on top of [hastebin](https://github.com/seejohnrun/haste-server).


It can be easily installed behind any network on any system that supports the following:

* A Node.js version which supports ES2015 Decorators.
* Can run one of the supported databases [<sup>[see]</sup>](#how-it-differs-from-hastebin).

#### **How it differs from hastebin:**

* Only supports phonetic key generation.
  * Although, this could be changed in the future and the source code is setup for a change like this.
* Written in TypeScript
* Doesn't support flat file storage.
* Supports far more databases.
  * Glue supports the following databases: **MySQL, MariaDB, Postgres, CockroachDB, SQLite, Microsoft SQL Server, Oracle, SAP Hana, sql.js, and MongoDB.**
* A well-documented RESTful API.
* Fairly well-documented source code.

## Installation

1. Download the package, and expand it
2. Explore the settings inside of `src/config.ts`, but the defaults should be good
3. Run `yarn` to install required packages
4. Run `yarn start`

`yarn start` automatically builds the source code. This should be a very fast process due to our use of the SWC compiler.

## Configuration Options

* `host (String)`: host to serve on
* `port (Number)`: port to serve on
* `keyLength (Number)`: length of keys to generate
* `maxDocumentLength (Number)`: max age of documents
* `staticMaxAge (Number)`: max age of static assets
* `useBrotli (Boolean)`: to use brotli or to not use brotli
* `useGzip (Boolean)`: to use gzip or to not use gzip
* [`dbOptions`](#database-options)
* [`rateLimits`](#rate-limiting)
  * `requests (Number)`
  * `every (Number)`

### Rate Limiting

We use `koa-ratelimit` as our rate-limiter.

Right now, we use a [map](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Map) for storing data requested to rate-limiting. `koa-ratelimit` supports an option for Redis which may be added to the config in the future.

The option `requests` is passed through to the `max` option in koa-ratelimit, and the `every` option is passed through to the `duration` option.

**NOTE: `every` is a millisecond value.**

### Database options

Section currently being written.

## Author

* Originally by [John Crepezzi <john.crepezzi@gmail.com>](https://github.com/seejohnrun)
* Rewritten by [Luke Whrit <me@lukewhrit.xyz>](https://github.com/324Luke)

## License

Glue is licensed under the permissive MIT license, same as [haste-server](https://github.com/seejohnrun/haste-server).

> **Copyright 2020 Luke Whrit**

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

> A copy of this license can also be found in the [`LICENSE.md`](LICENSE.md) file.
