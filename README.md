# Glue

[![Requires.io](https://img.shields.io/requires/github/324Luke/glue)](https://requires.io/github/324Luke/glue/requirements/?branch=master) [![GitHub last commit](https://img.shields.io/github/last-commit/324Luke/glue)](https://github.com/324Luke/glue/commits/master) [![TravisCI](https://img.shields.io/travis/324Luke/glue)](https://travis-ci.org/github/324Luke/glue)

Glue is an modern fork of [hastebin](https://github.com/seejohnrun/haste-server). It can easily be installed behind any network and on any system node.js supports.

Currently a public version does not exist, eventually the goal is to have one though.

**WARNING:** This readme may not reflect glue's current functionality

How it differs from hastebin:

* Always uses phonetic key generation
* Written in TypeScript
* Pure Node.js, no C++ extensions
* Supports more databases, by way of sequelize
* No storing in file system
* No support for Redis, Memcached, or RethinkDB

## Installation

1. Download the package, and expand it
2. Explore the settings inside of config.js, but the defaults should be good
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
  * Rate Limits
    * `requests (Number)`
    * `every (Number)`

## Rate Limiting

When present, the `rateLimits` option enables built-in rate limiting courtesy
of `connect-ratelimit`.  Any of the options supported by that library can be
used and set in `config.json`.

See the README for [connect-ratelimit](https://github.com/dharmafly/connect-ratelimit)
for more information!

## Storage

This section is currently being written.

## Author

* Originally by John Crepezzi <john.crepezzi@gmail.com>
* Rewrite by Luke Whrit <me@lukewhrit.xyz>

## License

(The MIT License)

Copyright © 2011-2012 John Crepezzi

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the ‘Software’), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED ‘AS IS’, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE

### Other components:

* jQuery: MIT/GPL license
* highlight.js: Copyright © 2006, Ivan Sagalaev
* highlightjs-coffeescript: WTFPL - Copyright © 2011, Dmytrii Nagirniak
