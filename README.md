# Glue

[![Requires.io](https://img.shields.io/requires/github/324Luke/glue)](https://requires.io/github/324Luke/glue/requirements/?branch=master) [![GitHub last commit](https://img.shields.io/github/last-commit/324Luke/glue)](https://github.com/324Luke/glue/commits/master) [![TravisCI](https://img.shields.io/travis/324Luke/glue)](https://travis-ci.org/github/324Luke/glue)

> **⚠️ WARNING:** This readme may not reflect glue's current functionality

Glue is a modern pastebin service, built on top of [`seejohnrun/hastebin-server`](https://github.com/seejohnrun/hastebin-server).

It can easily be installed on virtually any system as long as it supports a recently modern version of NodeJS (8+ should be fine) & a [compatible](#how-it-differs-from-hastebin) database.

A public hosted version will be available once glue hits a stable state.

#### **How it differs from hastebin:**

* Always uses phonetic key generation
* Written in TypeScript
* Always stores in database
* AlpineJS frontend
* Supports MySQL, PostgreSQL, Microsoft SQL Server, Oracle DB & SQLite 3

## Installation

1. Download the package, and expand it
2. Explore the settings inside of `src/config.ts`, but the defaults should be good
3. Run `yarn` to install required packages
4. Run `yarn start`

## Configuration Options

* Top Level
  * `host (String)`
  * `port (Number)`
  * `keyLength (Number)`
  * `maxLength (Number)`
  * `staticMaxAge (Number)`
  * `recompressStaticAssets (Boolean)`
  * [`dbOptions`](#database-options)
  * [`rateLimits`](#rate-limiting)
    * `requests (Number)`
    * `every (Number)`

### Rate Limiting

Section currently being written.

### Database options

Section currently being written.

## Author

* Originally by John Crepezzi <john.crepezzi@gmail.com>
* Rewrite by Luke Whrit <me@lukewhrit.xyz> (Discord: `Luke#1000`)

## License (MIT)

**Copyright 2020 Luke Whrit**

**Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:**

> **The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.**

> **THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.**
