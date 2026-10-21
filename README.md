<p align="center">
    <img
        width="600"
        src="https://raw.githubusercontent.com/lukewhrit/spacebin/main/.github/assets/spacebin-text-logo/Spacein%20Text.svg"
        alt="spacebin - hastebin fork focused on stability and maintainability"
    />
</p>

# 🚀 Spacebin

[![codecov](https://codecov.io/gh/lukewhrit/spacebin/graph/badge.svg?token=NNZDS74DB1)](https://codecov.io/gh/lukewhrit/spacebin) [![GitHub license](https://img.shields.io/github/license/lukewhrit/spacebin?color=%20%23e34b4a&logoColor=%23000000)](LICENSE) [![Build](https://github.com/lukewhrit/spacebin/actions/workflows/build.yml/badge.svg?branch=develop)](https://github.com/lukewhrit/spacebin/actions/workflows/build.yml) 

Spacebin is a modern Pastebin server implemented in Go and is capable of serving notes, novels, code, or any other form of text.

Spacebin was designed to be fast and reliable, avoiding the problems of many current pastebin servers. Spacebin highlights code on the server, so everything works with JavaScript disabled except encrypted pastes, which are encrypted in your browser. We have many more features in the works. It is entirely self-hostable, and available in a Docker image.

Pastebins are a type of online content storage service where users can store plain text document, e.g. program source code. For more information and the history of Pastebin see Wikipedia's [article on them](https://en.wikipedia.org/wiki/Pastebin).

> [!IMPORTANT]
> **Try our public online version at [https://spaceb.in](https://spaceb.in)**!

**Features:**

-   [x] 99% self-contained: only requires a database to run.
-   [x] Raw text and code uploading
-   [x] Phrase and random string identifiers.
-   [x] Custom named documents that are always available.
-   [x] Configurable ratelimiting, document size limits, ID formats, and more
-   [x] Modern user interface that works without JavaScript (only encrypted pastes need it)
-   [x] Syntax highlighting for all the most popular languages and Raw text mode
-   [x] SQLite, MySQL, and PostgreSQL Support
-   [x] Basic Auth for private instances
-   [X] Account system
-   [X] Reader view mode (Markdown is formatted and word wrapping is enabled)
-   [X] QR Codes
-   [X] [Password-protected encrypted pastes](#encrypted-pastes), encrypted in the browser so the server can't read them
-   [X] [Paste collections](#collections) for sharing several documents under one link

**Vote on future features: [Image/file uploading](https://github.com/lukewhrit/spacebin/discussions/446)**

## Table of Contents

- [🚀 Spacebin](#-spacebin)
  - [Table of Contents](#table-of-contents)
  - [Documentation](#documentation)
    - [Self-hosting](#self-hosting)
      - [Using Docker](#using-docker)
      - [Docker Compose](#docker-compose)
      - [Manually](#manually)
      - [Environment Variables](#environment-variables)
        - [Database Connection URI](#database-connection-uri)
    - [Usage](#usage)
      - [On the Web](#on-the-web)
        - [Encrypted Pastes](#encrypted-pastes)
        - [Collections](#collections)
      - [CLI](#cli)
    - [API](#api)
  - [Credits](#credits)
  - [Vulnerabilities](#vulnerabilities)
  - [License](#license)

## Documentation

> [!TIP]
> Need support? Visit our [Discussions](https://github.com/lukewhrit/spacebin/discussions) tab to get help from the community.

### Self-hosting

#### Using Docker

```sh
# Pull and run docker image on port 80
$ sudo docker pull spacebinorg/spirit
$ sudo docker run -d -e SPIRIT_CONNECTION_URI="sqlite://database.sqlite" -p 80:9000 spacebinorg/spirit
```

#### Docker Compose

Use the following config file to host Spacebin via Docker Compose:

```yml
services:
  spacebin:
    image: spacebinorg/spirit:latest
    restart: always
    environment:
      - SPIRIT_CONNECTION_URI=postgres://spacebin:password@postgres:5432/spacebin?sslmode=disable
    ports:
      - 9000:9000
    depends_on:
      - postgres
  postgres:
    image: postgres:16.3-alpine
    restart: always
    environment:
      - POSTGRES_USER=spacebin
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB=spacebin
    volumes:
      - postgres:/var/lib/postgresql/data

volumes:
  postgres:
```

#### Manually

> [!IMPORTANT]
> **Requires: [Git](https://git-scm.com/downloads), [Go 1.26.2](https://go.dev/doc/install), [GNU Makefile](https://www.gnu.org/software/make/#download), and a SQLite, MySQL, or [PostgreSQL](https://www.postgresql.org/download/) [server](https://m.do.co/c/beaf675c3e00).**

```sh
# Clone the Github repository
$ git clone https://github.com/lukewhrit/spacebin.git
$ cd spacebin

# Build the binary
$ make

# Start Spacebin
$ SPIRIT_CONNECTION_URI="sqlite://database.sqlite" ./bin/spacebin # SQLite
$ SPIRIT_CONNECTION_URI="mysql://<user>:<password>@tcp(<host>:3306)/<database>" ./bin/spacebin # MySQL
$ SPIRIT_CONNECTION_URI="postgres://<your PostgreSQL instance URI>" ./bin/spacebin # PostgreSQL

# Success! Spacebin is now available at port 9000 on your machine.
```

#### Environment Variables

| Variable Name                    | Type              | Default      | Description                                                                                                                      |
| -------------------------------- | ----------------- | ------------ | -------------------------------------------------------------------------------------------------------------------------------- |
| `SPIRIT_HOST`                    | String            | `0.0.0.0`    | Host address to listen on                                                                                                        |
| `SPIRIT_PORT`                    | Int               | `9000`       | HTTP port to listen on                                                                                                           |
| `SPIRIT_CONNECTION_URI`          | String            | **Required** | [Database connection URI](#database-connection-uri)                                                                              |
| `SPIRIT_RATELIMITER`             | String            | `200x5`      | Requests allowed across all clients combined, as `<requests>x<seconds>`: the default allows 200 every 5 seconds                  |
| `SPIRIT_HEADLESS`                | Bool              | `false`      | Set to `true` to disable the web interface and only serve the API                                                                |
| `SPIRIT_ANALYTICS`               | String            | `""`         | `<script>` tag for analytics (leave blank to disable). You'll also need to allow it in `SPIRIT_CSP`                              |
| `SPIRIT_CSP`                     | String            | See below    | `Content-Security-Policy` header sent with every response                                                                        |
| `SPIRIT_USERNAME`                | String            | `""`         | Basic Auth username. Basic Auth is enabled when both a username and password are set                                             |
| `SPIRIT_PASSWORD`                | String            | `""`         | Basic Auth password                                                                                                              |
| `SPIRIT_ACCOUNTS_ENABLED`        | Bool              | `false`      | Enables accounts, which let users sign up and edit or delete their documents and collections                                     |
| `SPIRIT_SESSION_TTL_HOURS`       | Int64             | `720`        | How long account sessions last, in hours                                                                                         |
| `SPIRIT_SESSION_COOKIE_SECURE`   | Bool              | `false`      | Only send the session cookie over HTTPS. Enable this if your instance uses HTTPS                                                 |
| `SPIRIT_SESSION_COOKIE_SAMESITE` | `lax` or `strict` | `lax`        | `SameSite` attribute of the session cookie                                                                                       |
| `SPIRIT_SESSION_COOKIE_DOMAIN`   | String            | `""`         | `Domain` attribute of the session cookie (leave blank for the current host)                                                      |
| `SPIRIT_ID_LENGTH`               | Int               | `8`          | Length for document and collection IDs                                                                                           |
| `SPIRIT_ID_TYPE`                 | `key` or `phrase` | `key`        | Format of IDs: `key` is a random string of letters and [`phrase` is a combination of words](https://github.com/lukewhrit/phrase) |
| `SPIRIT_MAX_SIZE`                | Int               | `400000`     | Max allowed length of a document, in characters                                                                                  |
| `SPIRIT_DOCUMENTS`               | []String          | `[]`         | List of any custom documents to serve                                                                                            |
| `SPIRIT_MAX_COLLECTION_SIZE`     | Int               | `25`         | Max number of documents in a [collection](#collections)                                                                          |

The default `SPIRIT_CSP` is:

```
default-src 'self'; frame-ancestors 'none'; base-uri 'none'; form-action 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline';
```

> [!NOTE]
> Encrypted pastes need the instance to be served over HTTPS (or opened on `localhost`), since browsers only allow encryption in secure contexts. Over plain HTTP the lock icon is hidden.

> [!WARNING]
> Environment variables for Spacebin are prefixed with `SPIRIT_`. They will be updated to `SPACEBIN_` in the next major version.

##### Database Connection URI

Spacebin supports three database formats: **SQLite**, **Postgres** and **MySQL**

-   For SQLite, use either the scheme `file://` or `sqlite://` and a file name.
    -   Example: `file://database.db`, or `sqlite:////var/lib/spacebin/database.db` for an absolute path
-   For PostgreSQL, use [the standard PostgreSQL URI format](https://stackoverflow.com/questions/3582552/what-is-the-format-for-the-postgresql-connection-string-url#20722229).
-   For MySQL, use the [DSN format](https://github.com/go-sql-driver/mysql?tab=readme-ov-file#dsn-data-source-name) prefixed with `mysql://` or `mariadb://`
    -   Example: `mysql://spacebin:password@tcp(localhost:3306)/spacebin`

### Usage

#### On the Web

To use Spacebin on the web, our team provides a web app. You can access the web app at **[spaceb.in](https://spaceb.in)**. You must use `https://spaceb.in/api` to access the API routes.

##### Encrypted Pastes

Click the lock icon next to the save button and choose a password before saving. Your document is encrypted in your browser before it's uploaded, so neither the server nor anyone with access to its database can read it.

When someone opens the link, they're asked for the password and the document is decrypted in their browser. Share the password separately from the link. **Passwords can't be recovered**, so if it's lost, so is the document.

Encrypted documents can't be edited, highlighted, or viewed in reader mode, as the server only has the encrypted version.

##### Collections

Collections share several documents under one link, with a tab for each document. Click the stack icon on the home page, then paste the IDs or links of the documents you want to include, one per line or separated by commas. You can also give the collection a title.

Collections are available at `/c/{collection}`. Anyone can create one. If you're signed in, it's listed on your account page, where you can delete it. Deleting a collection doesn't delete its documents.

#### CLI

Since Spirit supports `multipart/form-data` uploads, it's extremely easy to use on the command line via `curl`. The scripts also use `jq` so that you can get a machine-readable version of the document's ID, instead of a lengthy JSON object.

**To upload a string of text:**

```sh
curl -F content="Hello, world!" https://spaceb.in/api/ | jq -r .payload.id
```

**To upload from a file:**

```sh
curl -F content=@helloworld.txt https://spaceb.in/api/ | jq -r .payload.id
```

**To make a collection from existing documents:**

```sh
curl -F title="My notes" -F documents="WfwKGJfs, LcNmviWs" https://spaceb.in/api/collections | jq -r .payload.id
```

Encrypted pastes are encrypted by the client, so `curl` alone can't create or read them. See the [API](#api) for the format if you'd like to encrypt documents yourself.

### API

There are three primary API routes to: create a document, fetch a documents text content in JSON format, and fetch a documents **plain text** content. There are also two routes for creating and fetching collections of documents.

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
        "username": "",
        "encrypted": false,
        "created_at": "2023-08-06T00:01:33.143532-04:00",
        "updated_at": "2023-08-06T00:01:33.143532-04:00"
    }
}
```

-   Encrypted documents (see [Encrypted Pastes](#encrypted-pastes))
    -   Set `encrypted` to `true` (a boolean in JSON, the string `"true"` in forms) and send an encrypted envelope as the `content`
    -   Documents are encrypted in the browser, so the server never sees the password or the original text. The web app does this when you click the lock icon and choose a password
    -   The envelope is JSON: `{"v": 1, "iter": 600000, "salt": "…", "iv": "…", "ct": "…"}`. The key is derived from the password with PBKDF2-SHA256 using `iter` iterations (100,000 to 2,000,000) and the 16 byte `salt`, then the text is encrypted with AES-256-GCM using the 12 byte `iv`. `salt`, `iv` and `ct` (the ciphertext followed by the GCM tag) are standard base64
    -   The maximum document length applies to the envelope, which is about a third larger than the original text
    -   Encrypted documents can't be edited, and the fetch routes return the envelope with `"encrypted": true`

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
        "username": "",
        "encrypted": false,
        "created_at": "2023-08-06T00:01:33.143532-04:00",
        "updated_at": "2023-08-06T00:01:33.143532-04:00"
    }
}
```

-   `/api/{document}/raw`: Fetch Document - Raw
    -   `{document}` = Document ID
    -   Document ID lengths vary between instances. For `spaceb.in`, they will be exactly 8 characters
    -   Returns a `plain/text` file containing the content of the document.
-   `/api/collections`: Create Collection
    -   Groups existing documents under one link, viewable on the web at `/c/{collection}`
    -   Accepts JSON and multipart/form-data, and only accepts POST requests
    -   Include an optional `title` and a `documents` field listing document IDs or links to documents
        -   In JSON, `documents` is an array. In forms, IDs can be separated by commas or whitespace
    -   Instances are able to specify a maximum number of documents per collection
    -   Successful requests return a JSON body with the following format:

```json
{
    "error": "",
    "payload": {
        "id": "rQTJvpKB",
        "title": "My notes",
        "username": "",
        "documents": ["WfwKGJfs", "LcNmviWs"],
        "created_at": "2026-10-21T18:47:22Z"
    }
}
```

-   `/api/collections/{collection}`: Fetch Collection
    -   `{collection}` = Collection ID
    -   Returns the same JSON body as Create Collection. Documents that have been deleted are left out.

> [!TIP]
> There are two additional non-API routes: `/ping`: returns a 200 OK if the service is online, and `/config`: returns a JSON body with the instances configuration settings.

## Credits

Spacebin is a project designed and maintained by Luke Whritenour. Spacebin started out as a fork of [hastebin](https://github.com/toptal/haste-server). Although it no longer contains _any_ code from the original, we'd like to acknowledge our roots regardless.

Spacebin itself is built using Go and various libraries (i.e. [Chi](https://github.com/go-chi/chi), [pq](https://github.com/lib/pq), [Ozzo Validation](https://github.com/go-ozzo/ozzo-validation), [Chroma](https://github.com/alecthomas/chroma), [env](https://github.com/caarlos0/env)).

A full list of code contributors is available [on Github](https://github.com/lukewhrit/spacebin/graphs/contributors). We'd also like to thank [@jackdorland](https://github.com/jackdorland) for designing our logo/brand.

## Vulnerabilities

The Spacebin team takes security very seriously. If you detect a vulnerability please contact us via email (<hello@spaceb.in>) or create a vulnerability report via Github. We request that you hold off on publishing any vulnerabilities until after they've been patched, or at least 60 days have passed since you reported it. See our [security policy](.github/SECURITY.md) for full details.

## License

Spacebin is licensed under the Apache 2.0 license. A copy of this license can be found within the [`LICENSE`](LICENSE) file.
