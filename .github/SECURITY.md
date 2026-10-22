# Security Policy

The Spacebin team takes security very seriously. Thank you for helping keep Spacebin and its users safe.

## Supported Versions

Security fixes are released for the latest minor version only. If you're running an older version, please upgrade before reporting.

| Version | Supported          |
| ------- | ------------------ |
| 1.1.x   | :white_check_mark: |
| < 1.1   | :x:                |

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues, discussions, or pull requests.**

Instead, report them privately through either of these channels:

-   **Email (preferred):** [hello@spaceb.in](mailto:hello@spaceb.in)
-   **GitHub:** [open a private vulnerability report](https://github.com/lukewhrit/spacebin/security/advisories/new).

Please include as much of the following as you can:

-   The type of issue (e.g. XSS, SQL injection, denial of service)
-   The affected version, commit, or instance (e.g. `spaceb.in`)
-   The relevant configuration (database driver, environment variables) if it matters
-   Step-by-step instructions to reproduce the issue
-   A proof-of-concept or exploit code, if possible
-   The impact of the issue, including how an attacker might exploit it

## What to Expect

-   We'll acknowledge your report within **7 days**.
-   We'll confirm whether the issue is valid and keep you updated as we work on a fix.
-   Once a fix is released, we'll publish a GitHub Security Advisory and credit you for the discovery, unless you'd prefer to stay anonymous.

## Disclosure Policy

We ask that you hold off on publicly disclosing a vulnerability until it has been patched, or until **60 days** have passed since you reported it, whichever comes first.

## Scope

In scope:

-   The Spacebin server and API in this repository
-   The browser-side encryption used by encrypted pastes, including any way for the server or anyone else to read an encrypted paste without its password
-   The official Docker image built from this repository
-   The hosted instance at [spaceb.in](https://spaceb.in)

Out of scope:

-   Vulnerabilities in third-party dependencies that are already publicly known; please report those upstream
-   Self-hosted instances that are misconfigured, out of date, or modified
-   Content that users have uploaded to a Spacebin instance; report abusive content to that instance's operator
-   Volumetric denial-of-service attacks, spam, or social engineering
-   Findings from automated scanners without a demonstrated impact

## Safe Harbor

We won't pursue legal action against anyone who researches and reports vulnerabilities in good faith under this policy. Please:

-   Only test against your own self-hosted instance where possible
-   Avoid accessing, modifying, or deleting other users' documents
-   Avoid degrading service for other users of `spaceb.in`
