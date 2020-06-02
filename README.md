<p align="center">
  <img
    width="600"
    src="https://i-hate-the-windows.nt-kernel.design/2dYDKc8.png"
    alt="spacebin - open source hastebin fork focused on stability and maintainability"
  />
</p>
<p align="center">
  <a href="https://github.com/324Luke/spacebin/commits/master">
    <img
      src="https://img.shields.io/github/last-commit/324Luke/spacebin"
      alt="Latest Commit"
    />
  </a>
  <a href="https://requires.io/github/324Luke/Spacebin/requirements/?branch=master">
    <img
      src="https://img.shields.io/requires/github/324Luke/glue"
      alt="Requirements"
    />
  </a href="https://travis-ci.org/github/324Luke/spacebin">
    <img
      src="https://img.shields.io/travis/324Luke/spacebin"
      alt="Build Status"
    />
  </a>
  <a>
</p>

<p align="center">
<b>Spacebin is a modern pastebin service. Built ontop of John Crepezzi's hastebin, it focuses on stability and maintainability.</b></p>

* Stable and Maintainable, thanks to Typescript.
* Supports a [large amount of databases](#database-setup).
* A well-documented RESTful API.
* Easy to use and maintain.
* Highly customizable.

## Installation

1. Download the package, and expand it.
2. See the [Database Setup](#database-setup) section for database setup.
3. Explore the settings inside of `src/config.ts`, you'll most likely need to modify the `dbOptions` section but other than that you should be good to go with the defaults.
4. Run `yarn` to install required packages.
5. Run `yarn start`.
6. You'll now be able to see the service running on the port you configured, by default it's `7777`.

### Database Setup
**Spacebin will default to [PostgreSQL](https://www.postgresql.org) if no other database is specified**

First off, make sure you have a supported database. Spacebin uses TypeORM so you'll (most likely) have one.

Spacebin supports:
  * **MySQL**
  * **MariaDB**
  * **PostgreSQL** (default)
  * **CockroachDB**
  * **SQLite** (good fallback option if none others avaliable)
  * **Microsoft SQL Server**
  * **Oracle Database**
  * **SAP Hana**
  * **sql.js**
  * **MongoDB**

We recommend reading [TypeORM's documentation](https://typeorm.io/#/) on how to setup your particular database.

## Configuration Options

* `host (String)`: host to serve on
* `port (Number)`: port to serve on
* `keyLength (Number)`: length of keys to generate
* `maxDocumentLength (Number)`: max age of documents
* `staticMaxAge (Number)`: max age of static assets
* `useBrotli (Boolean)`: to use brotli or to not use brotli
* `useGzip (Boolean)`: to use gzip or to not use gzip (that is the question)
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

Spacebin was made possible by contributions from the Open Source community, as well as a few projects and people that stand out the most:

* Spacebin made by [Luke Whrit <me@lukewhrit.xyz>](https://github.com/324Luke)
* Hastebin originally by [John Crepezzi <john.crepezzi@gmail.com>](https://github.com/seejohnrun)
* Icon, graphic design and frontend contributions by [Jack Dorland <puggo@puggo.space>](https://github.com/heyitspuggo)
* Default color scheme design provided by [Jared Gorski's `spacecamp`](https://github.com/jaredgorski/spacecamp)
* And all the [other awesome contributors!](https://github.com/324Luke/spacebin/graphs/contributors).

## License

Spacebin is licensed under the very permissive "Clear BSD License".

This license allows for use in commercial & private situations and for distribution and modification of the source code.

Spacebin does **not** provide any warranty, does **not** hold any liability, and does **not** grant patent rights to contributors.

This license can also be found in markdown format in [LICENSE.md](LICENSE.md).

```
The Clear BSD License

Copyright (c) 2020 Luke Whrit
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted (subject to the limitations in the disclaimer
below) provided that the following conditions are met:

     * Redistributions of source code must retain the above copyright notice,
     this list of conditions and the following disclaimer.

     * Redistributions in binary form must reproduce the above copyright
     notice, this list of conditions and the following disclaimer in the
     documentation and/or other materials provided with the distribution.

     * Neither the name of the copyright holder nor the names of its
     contributors may be used to endorse or promote products derived from this
     software without specific prior written permission.

NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY
THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A
PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR
CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR
BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER
IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
```
