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

The one thing this section won't cover is installing [Git](https://git-scm.org), it is required but we assume you probsably already have that installed.

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

...

### Installing dependencies and your first build

...

## Styling of Code

* Use **whitespace** extensively to separate unrelated code statements. All code blocks should be surrounded by whitespace.
* Always try to follow the **existing style**.
* Prefer **readability and simplicity over efficiency**.
* Always run your code through **`go fmt`**.
* Attempt to **not** let any line be **longer than 80 characters**.

## Commit Guidelines

...
