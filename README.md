# Gitria
Gitria is a Git platform focused on keeping Git at its core.
It provides the essential Git hosting functionality without forcing users to adopt
additional features such as pull requests, issues, or project management.

These features are available, but they are completely optional.
Users and companies can choose whether they want to use them or integrate their existing tools instead.

# Why
Most Git hosting solutions try to provide an all-in-one experience, combining Git hosting with pull
requests, issues, project management, and other workflow features.

I don’t believe this approach should be the default.
Users and companies should be free to choose the tools and workflows that work best for them.

Gitria provides solid, straightforward Git integration without dictating how teams should work with their projects.
Features such as PRs and issues are there when you need them, but they never get in the way when you don’t.

The goal is simple: provide the Git platform, and let users decide how they want to work.

# Development
## Prerequisites
- Go 1.26 or later
- GNU Make
- OpenSSH
- Git

## Verify your Go installation:
```sh
go version
```

## Setup
Clone the repository and install the Go dependencies:
```sh
go mod download
```

Generate an SSH host key for local development:
```sh
ssh-keygen -t ed25519 -f ./id_rsa -N ""
```

This creates the test host key used by the development server.

## Run
Start the development server with:
```sh
make run
```

By default, the server listens on port 2222 and uses ./id_rsa as its SSH host key.

You can override the configuration when needed:
```sh
make run SSH_LISTEN_PORT=3333
```

## Build
Build the project:
```sh
make build
```
The binary is created at:
```sh
bin/gitria
```
Run it with:

```sh
./bin/gitria
```
