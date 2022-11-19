<p align="center">
    <img
        width="800"
        src="https://github.com/orca-group/wiki/blob/master/assets/spacebin-text-logo/github-banner.png?raw=true"
        alt="spacebin - hastebin fork focused on stability and maintainability"
    />
</p>

<p align="center">
    <a href="https://github.com/orca-group/spirit/commits/master">
        <img
            src="https://img.shields.io/github/last-commit/orca-group/spirit"
            alt="Latest Commit"
        />
    </a>
      <a href="https://github.com/orca-group/spirit/curiosity/master/LICENSE.md">
        <img
            alt="GitHub"
            src="https://img.shields.io/github/license/orca-group/spirit?color=%20%23e34b4a&logoColor=%23000000"
        />
    </a>
    <a href="https://app.codacy.com/gh/orca-group/spirit">
        <img
              alt="Codacy code quality grade"
              src="https://img.shields.io/codacy/grade/1aeee90534ef4dbeb4efa3fa6b4d027c"
        />
    </a>
    <a href="https://github.com/orca-group/spirit/workflows/build">
        <img
            alt="Build Status"
            src="https://github.com/orca-group/spirit/workflows/build/badge.svg"
        />
    </a>
    <a href="https://goreportcard.com/report/github.com/orca-group/spirit">
        <img
            alt="Go Report Card"
            src="https://goreportcard.com/badge/github.com/orca-group/spirit"
        />
    </a>
</p>

> **🚀 Spirit is the primary implementation of the Spacebin Server, written in Go and maintained by the Orca Group.**
>\
>\
> [**📖 Documentation**](https://docs.spaceb.in) | [**🚀 Live Instance (with Pulsar)**](https://spaceb.in)

## 🚀 What is Spacebin?

Spacebin is a highly-reliable pastebin server, built with Go, that's capable of serving notes, novels, code or any other form of text!

Pastebin's are a type of online content storage service where users can store plain text document, e.g. program source code. For more information and the history of Pastebin see Wikipedia's [article on them](https://en.wikipedia.org/wiki/Pastebin).

## ☄️ Clients

The Orca Group primarily maintains Spirit, but we also run a few clients for interacting with the server.

These are: [**🌟 Pulsar**](https://github.com/orca-group/pulsar)&mdash;a lightweight web client written in Svelte, and [**☄️ Comet**](https://github.com/orca-group/comet)&mdash;a speedy POSIX sh CLI for Spirit.

Our community has also contributed some great clients alternative to ours! You can check out a list of them [here, on the Github Topic](https://github.com/topics/spacebin).

## 🖨️ Self-hosting

**Requires: Docker**

```sh
# Pull and run docker image on port 80
$ sudo docker pull spacebinorg/spirit
$ sudo docker run -d -p 80:9000 spacebinorg/spirit
```

## 👥 Contributors

Spacebin (and Spirit) is a project by Luke Whritenour, associated with the [Orca Group](https://github.com/orca-group)&mdash;a developer collective. Spirit was forked from [hastebin](https://github.com/toptal/haste-server) by John Crepezzi (now managed by Toptal), and although it no longer contains **any** code from the original we'd like to thank him regardless. Spirit itself is built using [Fiber](https://gofiber.io/), [Gorm](https://gorm.io), [Ozzo Validation](https://github.com/go-ozzo/ozzo-validation), [Cron](https://github.com/robfig/cron), [Koanf](https://github.com/knadh/koanf), and (of course) [Go](https://go.dev/) itself!

You can see a full list of code contributors to Spirit [here, on Github](https://github.com/orca-group/spirit/graphs/contributors).

Additionally, we'd like to thank [@uwukairi](https://github.com/uwukairi) for designing our logo/brand, and [Brett Bender](https://github.com/GreatGodApollo) for additional technical help.

## 👮 Vulnerabilities

The Spacebin team takes security very seriously. If you detect a vulnerability please contact `lukewhrit@proton.me`. We request that you hold of on publishing any vulnerabilities until after they've been patched, or at least 60 days have passed since you reported it.

## 📑 License and Copyright

Spirit is licensed under the Apache 2.0 license. A copy of this license can be found within the [`LICENSE`](LICENSE) file.
