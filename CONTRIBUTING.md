# Contributing to Spacebin

Hey! First off: Thanks for deciding to contribute to Spirit, but, just before you get started we need to go over a few things.

When contributing to Spirit, please keep in mind we have a few fundamental guidelines that you must abide by:

* Attempt to create fast, elegant and readable code.
 * Beginners are welcome! Any PR you will submit will be subject to code review. If we deem it's not good enough, we will try to provide insight on how to make it better.
 * Readability is important, make sure to leave comments on your code where ever you can.
* Make thoughtful decisions about Pull Requests and Issues.
* Follow the [Code of Conduct](code_of_conduct.md).

## Setting up your environment

This section covers installing and configuring your development environment for contributing to Spirit.

The one thing this section won't cover is installing [Git](https://git-scm.org), it is required but we assume you probably already have that installed.

### Installing Golang

Since Spirit is written in Go, we're going to need it in our installed in our environments. Obviously, you can feel free to skip this part if already completed.

The steps for installing Go can drastically differ operating system to operating system; So, you should really be looking for instructions from the Go [documentation pages](https://golang.org/doc/install).

You'll need **at least** Go 1.14.

### Installing Magefile

We use Magefile as our build system — a make/rake-like build tool using Go. Installation is easy, just run these three commands:

```sh
$ git clone https://github.com/magefile/mage
$ cd mage
$ go run bootstrap.go
```

You can run these commands in any directory of your choice.

### Cloning the Git repository

In order to actually push modifications to a Spirit project, you need to first create a fork of it on your personal account or an organization in which you have permission to at minimum edit and create a repository.

Github has in-depth documentation on how to create a fork, which you can read [here](https://docs.github.com/en/github/getting-started-with-github/fork-a-repo). If you do not already know how to create a fork, you should read over this webpage before going any further.

After forking and cloning the repository, you need to change which branch your working on. Run the command `git checkout develop`, this will move your local clone to `develop` branch of Spirit.

### Installing dependencies and your first build

Installing dependencies is quite easily done with Go, especially so ever since the introduction of [Modules](https://blog.golang.org/using-go-modules) in Go 1.11. You can just run **`go get ./...`**, that will install all of the packages defined in Spirit's `go.mod` file.

#### Building and running the program

...

### Make your changes

Go, make your changes to Spirit. Just make sure to follow our style guidelines and conventions and to rebuild and test your program with each change.

Additionally, once your done making changes you need to commit your code. Committing your code is quite easily done, just run the following two commands:

```sh
$ git add .
$ git commit -m "<Commit Message>"
```

**But, before you do this make sure you read over our [Commit Guidelines](#commit-guidelines) in this file.**

### Creating a Pull Request and submitting your changes

In order for your changes to be included in the next release you need to create a pull request. Like with forking, Github also includes documentation on how to create a pull request, which can be viewed [here](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request-from-a-fork).

Generally though, the process goes like this:

* First, Above the file directory on the main page, click the "Pull request" button.
* Second, In the "base branch" drop-down menu, select the branch of the upstream Spirit repository to merge changes into. This needs to be the "develop" branch, otherwise your PR will not be merged.
* Third, In the "head branch" drop-down, select the "develop" branch

## Styling of Code

* Use **whitespace** extensively to separate unrelated code statements. All code blocks should be surrounded by whitespace.
* Always try to follow the **existing style**.
* Prefer **readability and simplicity over efficiency**.
* Always run your code through **`go fmt`**.
* Attempt to **not** let any line be **longer than 80 characters**.

## Commit Guidelines

...
