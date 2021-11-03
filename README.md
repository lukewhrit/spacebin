<p align="center">
    <img
        width="800"
        src="https://github.com/spacebin-org/wiki/blob/master/assets/spacebin-text-logo/github-banner.png?raw=true"
        alt="spacebin - hastebin fork focused on stability and maintainability"
    />
</p>

<p align="center">
    <a href="https://github.com/coral-dev/spirit/commits/master">
        <img
            src="https://img.shields.io/github/last-commit/coral-dev/spirit"
            alt="Latest Commit"
        />
    </a>
      <a href="https://github.com/coral-dev/spirit/curiosity/master/LICENSE.md">
        <img
            alt="GitHub"
            src="https://img.shields.io/github/license/coral-dev/spirit?color=%20%23e34b4a&logoColor=%23000000"
        />
    </a>
    <a href="https://app.codacy.com/gh/coral-dev/spirit">
        <img
              alt="Codacy code quality grade"
              src="https://img.shields.io/codacy/grade/d4ed85470f4045b4b1909cb896509915"
        />
    </a>
    <a href="https://github.com/coral-dev/spirit/workflows/build">
        <img
            alt="Build Status"
            src="https://github.com/coral-dev/spirit/workflows/build/badge.svg"
        />
    </a>
    <a href="https://goreportcard.com/report/github.com/coral-dev/spirit">
        <img
            alt="Go Report Card"
            src="https://goreportcard.com/badge/github.com/coral-dev/spirit"
        />
    </a>
</p>

> **🚀 Spirit is the primary Spacebin server implementation. It is written in Golang and maintained by the Spacebin team.**
>\
>\
> [**📖 Documentation**](https://docs.spaceb.in) | [**🌟 Development Branch**](https://github.com/coral-dev/spirit/tree/develop) | [**🚀 Live Instance (with Pulsar)**](https://spaceb.in)

## 🚀 What is Spacebin?

Spacebin is a highly-reliable pastebin server, built with Go, that's capable of serving notes, novels, code or any other form of text!

Pastebin's are a type of online content storage service where users can store plain text document, e.g. program source code.

For more information and the history of Pastebin see Wikipedia's [article on them](https://en.wikipedia.org/wiki/Pastebin).

## ☄️ Clients

Along with Spirit, we maintain a small number of other clients for interacting with Spacebin via either the web or your terminal.

These clients are: [🌟 Pulsar](https://github.com/coral-dev/pulsar) &mdash; a lightweight web client written in Svelte, and [☄️ Comet](https://github.com/coral-dev/comet) &mdash; a speedy command-line program for Spirit written in Go.

The community around Spacebin has also developed a larger number of clients, you can view a nearly complete list maintained by the Spacebin Team, [here on our documentation site](https://docs.spaceb.in/clients_and_libraries.html). 

## ✍️ Contributing

Spacebin uses a lot of technologies and follows a lot of guidelines, all of these are detailed in [`CONTRIBUTING.md`](CONTRIBUTING.md) along with basic environment setup information.

## 🖨️ Self-hosting

**Requires: Docker**

```sh
# Pull and run docker image on port 80
$ sudo docker pull spacebinorg/spirit:v0.1.6a
$ sudo docker run -d -p 80:9000 spacebinorg/spirit:v0.1.6a
```

## 👥 Contributors

Spirit development is lead by Ava Whrit, other team members, and various other contributors from the open source community.

Here's a list of notable contributors to Spirit:

* [Ava Whrit <avaaxcx@pm.me>](https://github.com/avaaxcx) - Lead developer and maintainer.
* [Brett Bender <brett@brettb.xyz>](https://github.com/greatgodapollo) - Contributor.

## 👮 Vulnerabilities

The Spacebin team takes security very seriously. If you detect a vulnerability please contact `avaaxcx@pm.me`. 

We ask that you do not publish any vulnerabilities after they have been patched or after a 30 day period since being reported has passed.

## 📑 License and Copyright

Spirit is licensed under the Apache 2.0 license. A copy of this license can be found within the [`license`](license.md) file.
