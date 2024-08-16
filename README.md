<p align="center">
    <img
        width="800"
        src="https://github.com/lukewhrit/spacebin/blob/master/.github/assets/spacebin-text-logo/github-banner.png?raw=true"
        alt="spacebin - hastebin fork focused on stability and maintainability"
    />
</p>

# 🚀 Spacebin

[![codecov](https://codecov.io/gh/lukewhrit/spacebin/graph/badge.svg?token=NNZDS74DB1)](https://codecov.io/gh/lukewhrit/spacebin) [![GitHub license](https://img.shields.io/github/license/lukewhrit/spacebin?color=%20%23e34b4a&logoColor=%23000000)](LICENSE) [![Build](https://github.com/lukewhrit/spacebin/actions/workflows/build.yml/badge.svg?branch=develop)](https://github.com/lukewhrit/spacebin/actions/workflows/build.yml)
[![Go report card](https://goreportcard.com/badge/github.com/lukewhrit/spacebin)](https://goreportcard.com/report/github.com/lukewhrit/spacebin)

Spacebin is a modern Pastebin server implemented in Go and maintained by Luke Whritenour. It is capable of serving notes, novels, code, or any other form of text! Spacebin was designed to be fast and reliable, avoiding the problems of many current pastebin servers. Spacebin features JavaScript-based text highlighting, but works completely fine with JS disabled. Besides text highlighting, we have many more features in the works. It is entirely self-hostable, and available in a Docker image.

Pastebins are a type of online content storage service where users can store plain text document, e.g. program source code. For more information and the history of Pastebin see Wikipedia's [article on them](https://en.wikipedia.org/wiki/Pastebin).

**Features:**

-   [x] 99% self-contained: only requires Postgres to run.
-   [x] Raw text and file uploading
-   [x] Phrase and random string identifiers.
-   [x] Custom documents that are always available.
-   [x] Configurable ratelimiting, expiration, compression, etc.
-   [x] Modern, JavaScript-free user interface
-   [x] Syntax highlighting for all the most popular languages and Raw text mode
-   [ ] Password-protected encrypted pastes
-   [ ] SQLite Support
-   [ ] Paste collections
-   [ ] Reader view mode (Markdown is formatted and word wrapping is enabled)
-   [ ] QR Codes
-   [ ] Image uploading

> [!TIP] > **Try our public online version at [https://spaceb.in](https://spaceb.in)**!

## Table of Contents

- [🚀 Spacebin](#-spacebin)
  - [Table of Contents](#table-of-contents)
  - [Documentation](#documentation)
    - [Self-hosting](#self-hosting)
      - [Using Docker](#using-docker)
      - [Manually](#manually)
      - [Environment Variables](#environment-variables)
    - [Usage](#usage)
      - [On the Web](#on-the-web)
      - [CLI](#cli)
    - [API](#api)
  - [Credits](#credits)
  - [Vulnerabilities](#vulnerabilities)
  - [License](#license)

## Documentation

### Self-hosting

#### Using Docker

```sh
# Pull and run docker image on port 80
$ sudo docker pull spacebinorg/spirit
$ sudo docker run -d -p 80:9000 spacebinorg/spirit
```

#### Manually

> [!IMPORTANT] > **Requires: [Git](https://git-scm.com/downloads), [Go 1.22.4](https://go.dev/doc/install), [GNU Makefile](https://www.gnu.org/software/make/#download), and a [PostgreSQL](https://www.postgresql.org/download/) [server](https://m.do.co/c/beaf675c3e00).**

```sh
# Clone the Github repository
$ git clone https://github.com/lukewhrit/spacebin.git
$ cd spacebin

# Build the binary
$ make spirit

# Start Spacebin
$ SPIRIT_CONNECTION_URI="<your PostgreSQL instance URI>" ./bin/spirit

# Success! Spacebin is now available at port 9000 on your machine.
```

#### Environment Variables

| Variable Name           | Type                  | Default      | Description                                                                                                                                        |
| ----------------------- | --------------------- | ------------ | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| `SPIRIT_HOST`           | String                | `0.0.0.0`    | Host address to listen on                                                                                                                          |
| `SPIRIT_PORT`           | Int                   | `9000`       | HTTP port to listen on                                                                                                                             |
| `SPIRIT_RATELIMITER`    | String                | `200x5`      | Requests allowed per second before the user is ratelimited                                                                                         |
| `SPIRIT_CONNECTION_URI` | String                | **Required** | [PostgreSQL Database URI String](https://stackoverflow.com/questions/3582552/what-is-the-format-for-the-postgresql-connection-string-url#20722229) |
| `SPIRIT_HEADLESS`       | Bool                  | `False`      | Enables/disables the web interface                                                                                                                 |
| `SPIRIT_ANALYTICS`      | String                | `""`         | `<script>` tag for analytics (leave blank to disable)                                                                                              |
| `SPIRIT_ID_LENGTH`      | Int                   | `8`          | Length for document IDs                                                                                                                            |
| `SPIRIT_ID_TYPE`        | `"key"` or `"phrase"` | `key`        | Format of IDs: `key` is a random string of letters and [`phrase` is a combination of words](https://github.com/lukewhrit/phrase)                   |
| `SPIRIT_MAX_SIZE`       | Int                   | `400000`     | Max allowed size of a document in bytes                                                                                                            |
| `SPIRIT_EXPIRATION_AGE` | Int64                 | `720`        | Amount of time to expire documents after                                                                                                           |
| `SPIRIT_DOCUMENTS`      | []String              | `[]`         | List of any custom documents to serve                                                                                                              |

> [!WARNING]
> Environment variables for Spacebin are prefixed with `SPIRIT_`. They will be updated to `SPACEBIN_` in the next major version.

### Usage

#### On the Web

To use Spacebin on the web, our team provides a web app. You can access the web app at **[spaceb.in](https://spaceb.in)**. You must use `https://spaceb.in/api` to access the API routes.

A version of spacebin that is built directly from the `develop` branch is also available at \*\*[staging.spaceb.in](https://staging.spaceb.in)

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

-   `/api/`: Create Document
    -   Accepts JSON and multipart/form-data
    -   For both formats, include document content in a `content` field
    -   Only accepts POST requests
    -   Instances are able to specify a maximum document length.
        -   `spaceb.in` uses a 4MB maximum size.
    -   Successful requests return a JSON body with the following format:

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

-   `/api/{document}`: Fetch Document
    -   `{document}` = Document ID
    -   Document ID lengths vary between instances. For `spaceb.in`, they will be exactly **8** characters.
    -   Upon successful request, returns a JSON body with the following format:

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

-   `/api/{document}/raw`: Fetch Document - Raw
    -   `{document}` = Document ID
    -   Document ID lengths vary between instances. For `spaceb.in`, they will be exactly 8 characters
    -   Returns a `plain/text` file containing the content of the document.

> [!TIP]
> There are two additional non-API routes that are documented: `/ping`: returns a 200 OK if the service is online, and `/config`: returns a JSON body with the instances configuration settings.

## Credits

Spacebin is a project designed and maintained by Luke Whritenour. Spacebin started out as a fork of [hastebin](https://github.com/toptal/haste-server). Although it no longer contains _any_ code from the original, we'd like to acknowledge our roots regardless.

Spacebin itself is built using Go and various libraries (i.e. [Chi](https://github.com/go-chi/chi), [pq](https://github.com/lib/pq), [Ozzo Validation](https://github.com/go-ozzo/ozzo-validation), [Cron](https://github.com/robfig/cron), [env](https://github.com/caarlos0/env)).

A full list of code contributors is available [on Github](https://github.com/lukewhrit/spacebin/graphs/contributors). We'd also like to thank [@jackdorland](https://github.com/jackdorland) for designing our logo/brand.

## Vulnerabilities

The Spacebin team takes security very seriously. If you detect a vulnerability please contact us: <hello@spaceb.in>. We request that you hold of on publishing any vulnerabilities until after they've been patched, or at least 60 days have passed since you reported it.

## License

Spacebin is licensed under the Apache 2.0 license. A copy of this license can be found within the [`LICENSE`](LICENSE) file.
