<p align="center">
    <img
        width="800"
        src="https://github.com/orca-group/wiki/blob/master/assets/spacebin-text-logo/github-banner.png?raw=true"
        alt="spacebin - hastebin fork focused on stability and maintainability"
    />
</p>

# 🚀 Spirit

[![codecov](https://codecov.io/gh/orca-group/spirit/branch/develop/graph/badge.svg?token=NNZDS74DB1)](https://codecov.io/gh/orca-group/spirit) [![GitHub license](https://img.shields.io/github/license/orca-group/spirit?color=%20%23e34b4a&logoColor=%23000000)](LICENSE) [![Build](https://github.com/orca-group/spirit/actions/workflows/build.yml/badge.svg?branch=develop)](https://github.com/orca-group/spirit/actions/workflows/build.yml)
[![Go report card](https://goreportcard.com/badge/github.com/orca-group/spirit)](https://goreportcard.com/report/github.com/orca-group/spirit)

Spirit is the primary implementation of the Spacebin Server, written in Go and maintained by the Orca Group. Spacebin itself is a standardized pastebin server, that's capable of serving notes, novels, code or any other form of text!

Pastebin's are a type of online content storage service where users can store plain text document, e.g. program source code. For more information and the history of Pastebin see Wikipedia's [article on them](https://en.wikipedia.org/wiki/Pastebin).

## Table of Contents

- [🚀 Spirit](#-spirit)
  - [Table of Contents](#table-of-contents)
  - [Documentation](#documentation)
    - [Self-hosting](#self-hosting)
    - [Usage](#usage)
      - [On the Web](#on-the-web)
      - [CLI](#cli)
    - [API](#api)
  - [Credits](#credits)
  - [Vulnerabilities](#vulnerabilities)
  - [License](#license)

## Documentation

### Self-hosting

**Using Docker**

```sh
# Pull and run docker image on port 80
$ sudo docker pull spacebinorg/spirit
$ sudo docker run -d -p 80:9000 spacebinorg/spirit
```

**Manually**

WIP

### Usage

#### On the Web

To use Spacebin on the web, our team provides a web app written in [Svelte](https://svelte.dev): [Pulsar](https://github.com/orca-group/pulsar). A public instance of Spacebin using this client is available at [https://spaceb.in](https://spaceb.in) (the `/api` route can be used to access Spirit itself).

#### CLI

Since Spirit supports `multipart/form-data` uploads, it's extremely easy to use on the command line via `curl`. The scripts also use `jq` so that you can get a machine-readable version of the document's ID, instead of a lengthy JSON object.

**To upload a string of text:**

```sh
curl -v -F content="Hello, world!" https://spaceb.in/api | jq payload.id
```

**To upload from a file:**

```sh
curl -v -F content=@helloworld.txt https://spaceb.in/api | jq payload.id
```

### API

Work in progress. Check out the documentation website: [docs.spaceb.in](https://docs.spaceb.in).

## Credits

Spacebin (and Spirit) is a project by Luke Whritenour, associated with the [Orca Group](https://github.com/orca-group)&mdash;a developer collective. Spirit was forked from [hastebin](https://github.com/toptal/haste-server) by John Crepezzi (now managed by Toptal), and although it no longer contains **any** code from the original we'd like to thank him regardless. Spirit itself is built using [Chi](https://github.com/go-chi/chi), and [pq](https://github.com/lib/pq), [Ozzo Validation](https://github.com/go-ozzo/ozzo-validation), [Cron](https://github.com/robfig/cron), [env](https://github.com/caarlos0/env), and (of course) [Go](https://go.dev/) itself!

You can see a full list of code contributors to Spirit [here, on Github](https://github.com/orca-group/spirit/graphs/contributors).

Additionally, we'd like to thank [@uwukairi](https://github.com/uwukairi) for designing our logo/brand.

## Vulnerabilities

The Spacebin team takes security very seriously. If you detect a vulnerability please contact us: <hello@spaceb.in>. We request that you hold of on publishing any vulnerabilities until after they've been patched, or at least 60 days have passed since you reported it.

## License

Spirit is licensed under the Apache 2.0 license. A copy of this license can be found within the [`LICENSE`](LICENSE) file.
