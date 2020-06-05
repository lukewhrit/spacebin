# Contributing

Hello! Thank you for deciding to contribute to spacebin 🚀🌌, but just before you get started we need to go over a few things.

## Guidelines

When contributing to spacebin, please keep in mind we have a few fundamental guidelines that you must abide by:

- Attempt to create speedy, readable and elegant code.
  - Beginners are welcome! Any PR you will submit will be subject to code review. If we deem it's not good enough, we will try to provide insight on how to make it better.
  - Readability is important, make sure to leave comments on your code where ever you can.
- Make thoughtful decisions about Pull Requests and Issues.
- Follow the [Code of Conduct](https://github.com/spacebin-for-astronauts/.github/blob/master/CODE_OF_CONDUCT.md).

## Specifications

- Each route (endpoint) has it's own controller `src/controllers` which is were functionality is implemented.
- Once a route has a controller, you can add an interface for accessing that controller as an API endpoint.
- Each file in `src/routes` is a section where new routes are implemented. They should follow a similar format to the Document route.

## Formatting

Spacebin follows the [Standard JS](https://standardjs.com) code style, and uses ESLint to enforce it.

So, when contributing make sure to follow this style, and before committing any changes run `yarn lint`.

## Commit Guidelines

- Spacebin follows the Conventional Commit standard.
- Spacebin also has commitizen setup.

When you're committing, it is recommended to not use `git commit`, and to instead use `npx git-cz`. This provides a simpler interface for creating commit messages that follow the Conventional Commit standard.

