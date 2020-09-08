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
</p>

> **🚀 Spirit is the primary Spacebin server implementation. It is written in Golang and maintained by the Spacebin team.**
>\
>\
> [**📖 Documentation**](https://docs.spaceb.in) | [**🌟 Development Branch**](https://github.com/spacebin-org/spirit/tree/develop) | [**🚀 More Information**](https://github.com/spacebin-org/spacebin#readme)

## Contributing

Spacebin uses a lot of technologies & follows a lot of rules, all of these are detailed in [`CONTRIBUTING.md`](CONTRIBUTING.md) along with basic environment setup information.

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

Spacebin development is lead by Luke Whrit and [the other team members](https://github.com/orgs/spacebin-org/teams/sever-team).

* [Luke Whrit <lukewhrit@gmail.com>](https://github.com/lukewhrit) - Lead developer and maintainer.

## License

Spacebin is licensed under the Apache 2.0 license. A copy of this license can be found within the [`license`](license.md) file.
