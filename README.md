<p align="center">
    <img
        width="800"
        src="https://github.com/lukewhrit/spacebin/blob/master/.github/assets/spacebin-text-logo/github-banner.png?raw=true"
        alt="spacebin - hastebin fork focused on stability and maintainability"
    />
</p>

# 🚀 Spacebin

[![codecov](https://codecov.io/gh/lukewhrit/spacebin/graph/badge.svg?token=NNZDS74DB1)](https://codecov.io/gh/lukewhrit/spacebin) [![GitHub license](https://img.shields.io/github/license/lukewhrit/spacebin?color=%20%23e34b4a&logoColor=%23000000)](LICENSE) [![Build](https://github.com/lukewhrit/spacebin/actions/workflows/build.yml/badge.svg?branch=develop)](https://github.com/lukewhrit/spacebin/spirit/actions/workflows/build.yml)
[![Go report card](https://goreportcard.com/badge/github.com/lukewhrit/spacebin)](https://goreportcard.com/report/github.com/lukewhrit/spacebin)

Spacebin is a modern Pastebin server implemented in Go and maintained by Luke Whritenour. It is capable of serving notes, novels, code, or any other form of text! Spacebin was designed to be fast and reliable, avoiding the problems of many current pastebin servers. Spacebin features JavaScript-based text highlighting, but works completely fine with JS disabled. Besides text highlighting, we have many more features in the works. It is entirely self-hostable, and available in a Docker image.

Pastebins are a type of online content storage service where users can store plain text document, e.g. program source code. For more information and the history of Pastebin see Wikipedia's [article on them](https://en.wikipedia.org/wiki/Pastebin).

**Features:**
- [X] 99% self-contained: only requires Postgres to run.
- [X] Raw text and file uploading
- [X] Phrase and random string identifiers.
- [X] Custom documents that are always available.
- [X] Configurable ratelimiting, expiration, compression, etc.
- [X] Modern, JavaScript-free user interface
- [X] Syntax highlighting for all the most popular languages and Raw text mode
- [ ] Password-protected encrypted pastes
- [ ] SQLite Support
- [ ] Paste collections
- [ ] Reader view mode (Markdown is formatted and word wraping is enbaled)
- [ ] QR Codes
- [ ] Image uploading

## Table of Contents

- [🚀 Spacebin](#-spacebin)
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

To use Spacebin on the web, our team provides a web app. You can access the web app at [spaceb.in](https://spaceb.in). You must use `https://spaceb.in/api` to access the API routes.

#### CLI

Since Spirit supports `multipart/form-data` uploads, it's extremely easy to use on the command line via `curl`. The scripts also use `jq` so that you can get a machine-readable version of the document's ID, instead of a lengthy JSON object.

**To upload a string of text:**

```sh
curl -v -F content="Hello, world!" https://spaceb.in/ | jq payload.id
```

**To upload from a file:**

```sh
curl -v -F content=@helloworld.txt https://spaceb.in/ | jq payload.id
```

### API

There are three primary API routes to: create a document, fetch a documents text content in JSON format, and fetch a documents **plain text** content.

* `/api/`: Create Document
  * Accepts JSON and multipart/form-data
  * For both formats, include document content in a `content` field
  * Only accepts POST requests
  * Instances are able to specify a maximum document length.
    * `spaceb.in` uses a 4MB maximum size.
  * Successful requests return a JSON body with the following format:
```json
{
	"error": "",
	"payload": {
		"id": "WfwKGJfs",
		"content": "hello",
		"created_at": "2023-08-06T00:01:33.143532-04:00",
		"updated_at": "2023-08-06T00:01:33.143532-04:00"
	}
}
```
* `/api/{document}`: Fetch Document
  * `{document}` = Document ID
  * Document ID lengths vary between instances. For `spaceb.in`, they will be exactly characters.
  * 
  * Upon successful request, returns a JSON body with the following format:
```json
{
	"error": "",
	"payload": {
		"id": "WfwKGJfs",
		"content": "hello",
		"created_at": "2023-08-06T00:01:33.143532-04:00",
		"updated_at": "2023-08-06T00:01:33.143532-04:00"
	}
}
```
* `/api/{document}/raw`: Fetch Document - Raw
  * `{document}` = Document ID
  * Document ID lengths vary between instances. For `spaceb.in`, they will be exactly 8 characters
  * Returns a `plain/text` file containing the content of the document.

There are two additional non-API routes that are documented: `/ping`: returns a 200 OK if the service is online, and `/config`: returns a JSON body with the instances configuration settings.

## Credits

Spacebin (and all other associated programs) is a project designed and maintained by Luke Whritenour. Spirit was forked from [hastebin](https://github.com/toptal/haste-server) by John Crepezzi (now managed by Toptal), and although it no longer contains **any** code from the original we'd like to thank him regardless. Spirit itself is built using [Chi](https://github.com/go-chi/chi), and [pq](https://github.com/lib/pq), [Ozzo Validation](https://github.com/go-ozzo/ozzo-validation), [Cron](https://github.com/robfig/cron), [env](https://github.com/caarlos0/env), and (of course) [Go](https://go.dev/) itself!

You can see a full list of code contributors to Spirit [here, on Github](https://github.com/lukewhrit/spacebin/graphs/contributors).

Additionally, we'd like to thank [@jackdorland](https://github.com/jackdorland) for designing our logo/brand.

## Vulnerabilities

The Spacebin team takes security very seriously. If you detect a vulnerability please contact us: <hello@spaceb.in>. We request that you hold of on publishing any vulnerabilities until after they've been patched, or at least 60 days have passed since you reported it.

## License

Spirit is licensed under the Apache 2.0 license. A copy of this license can be found within the [`LICENSE`](LICENSE) file.
