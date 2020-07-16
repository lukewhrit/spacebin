# Contributing

Hello! Thank you for deciding to contribute to spacebin 🚀🌌, but just before you get started we need to go over a few things.

Contributing to Spacebin requires sufficient knowledge in the following technologies:

- TypeScript
- Express
- TypeORM
- Jest & Supertest

And, of course, a Node.JS development environment.

## Guidelines

When contributing to spacebin, please keep in mind we have a few fundamental guidelines that you must abide by:

- Attempt to create fast, elegant and readable code.
  - Beginners are welcome! Any PR you will submit will be subject to code review. If we deem it's not good enough, we will try to provide insight on how to make it better.
  - Readability is important, make sure to leave comments on your code where ever you can.
- Make thoughtful decisions about Pull Requests and Issues.
- Follow the [Code of Conduct](CODE_OF_CONDUCT.md).

### Specifications

Spacebin follows a MVC like pattern, ERC (Entity, Route, Controller):

- **Entities** are TypeORM entities. They are how you interact with the database.
- **Routes** are Express router instances.
- **Controllers** are TypeScript classes that provide functionality to your routes and interact with your entities. They are the glue for your new feature.

#### Directory Structure

- `src/` - Where the source code for Spacebin is kept.
  - `__tests__/` - Jest test files.
  - `controllers/` - Classes that define the functionality of router endpoints.
  - `routes/` - Router instances which make functionality defined in `controllers` available to consumers.
  - `entities/` - Where our TypeORM entities are kept. See the [docs for more information on entities](https://typeorm.io/#/entities).
  - `validators/` - Basic objects that define the inputs and outputs of routes. See the format of `document.validator.ts` to learn how to write these.

All files must be suffixed with the singular form of a directories name, e.g. for a directory named entities the suffix would be `entity`.

### Code Formatting

Spacebin primarily follows the [JavaScript Standard Style](https://standardjs.com/).

Notably this includes rules such as:

- No semicolons
- Only semi-quotes

We also follow the recommended styles for TypeScript. All of this is enforced by ESLint.

Some other things:

- Try to follow the existing style
- Prefer readability over efficiency.

### Commit Guidelines

Before committing you must make sure your git client of choice supports [Git Hooks](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks).

Git hooks do two things for Spacebin:

1. Enforce the StandardJS code style via `lint-staged`.
2. Enforce consistent commit messages via `commitlint`. We follow the Conventional Commit standard.

**If you think this file is missing something please contact `Luke#1000` in [the Discord server](https://discord.gg/zsxwgYc).**
