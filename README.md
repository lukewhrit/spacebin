<!-- Spacebin Curiosity README.md -->

<p align="center">
  <img
    width="800"
    src="https://github.com/spacebin-org/assets/blob/master/assets/images/spacebin/icons-large/spacebin-large.png?raw=true"
    alt="spacebin - hastebin fork focused on stability and maintainability"
  />
</p>

<p align="center">
  <a href="https://github.com/spacebin-org/curiosity/commits/master">
    <img
      src="https://img.shields.io/github/last-commit/spacebin-org/curiosity?style=flat-square"
      alt="Latest Commit"
    />
  </a>
  <a href="https://discord.gg/hXxBtMJ">
    <img
      alt="Discord"
      src="https://img.shields.io/discord/717911514593493012?color=7289da&style=flat-square"
    />
  </a>
  <a href="https://github.com/spacebin-org/spirit/curiosity/master/LICENSE.md">
    <img
      alt="GitHub"
      src="https://img.shields.io/github/license/spacebin-org/curiosity?color=%20%23e34b4a&logoColor=%23000000&style=flat-square"
    />
  </a>
  <a href="https://app.codacy.com/gh/spacebin-org/curiosity">
    <img
      alt="Codacy code quality grade"
      src="https://app.codacy.com/project/badge/Grade/b15352aa5c394722948e4fc081ed1f60?style=flat-square"
    />
  </a>
</p>

> **🚀 Curiosity is the future of Spacebin's main server. It is written in Golang and is maintained by the Spacebin team.**
>\
>\
> [**📖 Documentation**](https://docs.spaceb.in) | [**🌟 Development Branch**](https://github.com/spacebin-org/curiosity/tree/develop) | [**🚀 More Information**](https://github.com/spacebin-org/spacebin#readme)

**Note: This will eventually be merged back into the Spirit repo, and subsequently be renamed to Spirit. If a new official server implementation were to arise it would be named Curiosity as well.**

## Contributing

Spacebin uses a lot of technologies & follows a lot of rules, all of these are detailed in [`CONTRIBUTING.md`](CONTRIBUTING.md) along with basic environment setup information.

## Self-hosting

**Requires: Git & Docker**

```sh
# Clone repository from git remote
$ git clone https://github.com/spacebin-org/curiosity.git
$ cd curiosity

# Do any configuration you may need to do now

# Build and run docker image on port 80
$ sudo docker build -t spacebin-curiosity .
$ sudo docker run -d -p 80:9000 spacebin-curiosity
```

## Contributors

Spacebin development is lead by Luke Whrit and [the other team members](https://github.com/orgs/spacebin-org/teams/sever-team).

* [Luke Whrit <lukewhrit@gmail.com>](https://github.com/lukewhrit) - Lead developer and maintainer.

## License

Spacebin is licensed under the GNU General Public License v3. A copy of this license can be found in markdown format in [LICENSE.md](LICENSE.md).
