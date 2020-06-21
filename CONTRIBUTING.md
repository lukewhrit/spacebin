# Contributing

Hello! Thank you for deciding to contribute to spacebin 🚀🌌, but just before you get started we need to go over a few things.

Contributing to Spacebin requires sufficient knowledge in the following technologies:

- TypeScript
- KoaJS and the following modules:
  - `koa-joi-router`
  - `koa-bodyParser`
  - We also use the other modules, however you probably will not need to work with them:
    - `koa/cors`
    - `koa-helmet`
    - `koa-morgan`
    - `koa-ratelimit`
- TypeORM
- AvaJS

And, of course, a Node.JS development environment.

## Guidelines

When contributing to spacebin, please keep in mind we have a few fundamental guidelines that you must abide by:

- Attempt to create fast, elegant and readable code.
  - Beginners are welcome! Any PR you will submit will be subject to code review. If we deem it's not good enough, we will try to provide insight on how to make it better.
  - Readability is important, make sure to leave comments on your code where ever you can.
- Make thoughtful decisions about Pull Requests and Issues.
- Follow the [Code of Conduct](CODE_OF_CONDUCT.md).

### Specifications

Spacebin follows a ERC (Entity, Route, Controller) model:

- **Entities** are TypeORM entities. They are how you interact with the database
- **Routes** are KoaJS (more specifically `koa-joi-router`) router functions.
- **Controllers** are TypeScript classes that provide functionality to your routes and interact with your entities. They are the glue for your new feature.

#### Directories

- `src/` - Where the source code for Spacebin is kept.
  - `controllers/` - Where the controllers we just discussed are stored.
  - `routes/` - This is where the files for the routes we talked about are kept.
  - `entities/` - Where our TypeORM entities are kept. See the [docs for more information on entities](https://typeorm.io/#/entities).
  - `validators/` - We use Joi for route validation, and it's objects can become very long. Therefore we opted to split them up into a separate file. These files are just plain JavaScript objects.
  - `tests/` - Every piece of functionality should be tested, this way we can ensure Spacebin is stable. We use `AvaJS` as our testing framework, see [their docs for information on how to write these files](https://github.com/avajs/ava/tree/master/docs). The suffix for test files is different from all others, test files must be suffixed with `spec`, instead of `test`.

All files must be suffixed with the singular form of a directories name, e.g. for a directory named entities the suffix would be `entity`.

### Code Formatting

Spacebin primarily follows the [JavaScript Standard Style](https://standardjs.com/).

Notably this includes rules such as:

- No semicolons
- Only semi-quotes

We also follow the recommended styles for TypeScript. All of this is enforced by ESLint, so please make sure your editor supports ESLint and before committing changes make sure to run `yarn lint` (Note: If your git client supports [hooks](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks) it will automatically do this for you). 

### Commit Guidelines

Before committing you must make sure your git client of choice supports [Git Hooks](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks).

Git hooks do two things for Spacebin:

1. Enforce the StandardJS code style via `lint-staged`.
2. Enforce consistent commit messages via `commitlint`. We follow the Conventional Commit standard.


**If you think this file is missing something please contact `Luke#1000` in [the Discord server](https://discord.gg/zsxwgYc).**
