<!-- Spacebin Curiosity README.md -->

<p align="center">
    <img
        width="800"
        src="https://github.com/spacebin-org/spacebin/blob/master/assets/images/spacebin/icons-large/spacebin-large.png?raw=true"
        alt="spacebin - hastebin fork focused on stability and maintainability"
    />
</p>

<p align="center">
    <a href="https://github.com/spacebin-org/spirit/commits/master">
        <img
            src="https://img.shields.io/github/last-commit/spacebin-org/spirit"
            alt="Latest Commit"
        />
    </a>
    <a href="https://discord.gg/hXxBtMJ">
        <img
            alt="Discord"
            src="https://img.shields.io/discord/717911514593493012?color=7289da"
        />
    </a>
      <a href="https://github.com/spacebin-org/spirit/curiosity/master/LICENSE.md">
        <img
            alt="GitHub"
            src="https://img.shields.io/github/license/spacebin-org/spirit?color=%20%23e34b4a&logoColor=%23000000"
        />
    </a>
    <a href="https://app.codacy.com/gh/spacebin-org/spirit">
        <img
              alt="Codacy code quality grade"
              src="https://img.shields.io/codacy/grade/ea24e2f7bf7d493e87a38cdcce4060b5"
        />
    </a>
    <a href="https://github.com/spacebin-org/spirit/workflows/build">
        <img
            alt="Build Status"
            src="https://github.com/spacebin-org/spirit/workflows/build/badge.svg"
        />
    </a>
    <a href="https://goreportcard.com/report/github.com/spacebin-org/spirit">
        <img
            alt="Go Report Card"
            src="https://goreportcard.com/badge/github.com/spacebin-org/spirit"
        />
    </a>
</p>

> **🚀 Spirit is the primary Spacebin server implementation. It is written in Golang and maintained by the Spacebin team.**
>\
>\
> [**📖 Documentation**](https://docs.spaceb.in) | [**🌟 Development Branch**](https://github.com/spacebin-org/spirit/tree/develop) | [**🚀 More Information**](https://github.com/spacebin-org/spacebin#readme)

## What is Spacebin?

Spacebin is a highly-reliable pastebin server, built with Go, that's capable of serving notes, novels, code or any other form of text!

Pastebin's are a type of online content storage service where users can store plain text document, e.g. program source code.

For more information and the history of Pastebin see Wikipedia's [article on them](https://en.wikipedia.org/wiki/Pastebin).

## Contributing

Spacebin uses a lot of technologies and follows a lot of guidelines, all of these are detailed in [`CONTRIBUTING.md`](CONTRIBUTING.md) along with basic environment setup information.

## Self-hosting

**Requires: Git & Docker**

```sh
# Clone repository from git remote
$ git clone https://github.com/spacebin-org/spirit.git
$ cd spirit
$ git checkout -b develop

# Build and run docker image on port 80
$ sudo docker build -t spacebin-server .
$ sudo docker run -d -p 80:9000 spacebin-server
```

## Contributors

Spirit development is lead by Luke Whrit, [other team members](https://github.com/orgs/spacebin-org/teams/sever-team), and various other contributors from the open source community.

Here's a list of notable contributors to Spirit:

* [Luke Whrit <lukewhrit@gmail.com>](https://github.com/lukewhrit) - Lead developer and maintainer.
* [Brett Bender <brett@brettb.xyz>](https://github.com/greatgodapollo) - Developed our testing infrastructure.

## Vulnerabilities

The Spacebin team takes security very seriously. If you detect a vulnerability please contact `lukewhrit@pm.me`. 

We ask that you do not publish any vulnerabilities after they have been patched or after a 30 day period since being reported has passed.

## License

Spirit is licensed under the Apache 2.0 license. A copy of this license can be found within the [`license`](license.md) file.
