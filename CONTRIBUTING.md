This document describes the full process of setting up a fully working local development environment and submitting your first contribution.

## Clone Repositories

Clone `frizzante`.

```go
git clone https://github.com/razshare/frizzante
```

Clone `frizzante-starter`.

```go
git clone https://github.com/razshare/frizzante-starter
```

> [!TIP]
> If you don't have direct access to these repositories you will need to fork your own `frizzante` and `frizzante-starter` repositories and clone those instead.
> 
> Then when you're done with your changes you will need to submit a pull request.

## Install Frizzante Cli

Install the frizzante cli.

You install it using the go cli

```go
go install github.com/razshare/frizzante@latest
```

Or you can [download the binaries directly from GitHub](https://github.com/razshare/frizzante/releases).

One way or another, make sure the frizzante cli is on your path.

If you're installing it using the go cli, make sure that go binaries are on your path

```sh
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

If you're using a precompiled binary, make sure that's on your path

```sh
export PATH=$PATH:/path/to/frizzante
```

## Configure Projects

Navigate to your `frizzante` local repository and configure the project with make.

```sh
make configure
```

This will install `bun` and `air` in `.gen`.

Navigate to your `frizzante-starter` local repository and configure the project with make.

```sh
make configure
```

This will install `bun` and `air` in `.gen`.

## Local Frizzante Package

You most likely will want to try out your `frizzante` changes locally.

By default, your local `frizzante-starter` will pull the `frizzante` package from the public remote repository.\
You can link your local `frizzante-starter` to your local version of `frizzante` instead, in order to have immediate feedback.

Navigate to your `frizzante-starter` repository and modify your `go.mod` file to replace the remote `frizzante` dependency with the local one using `replace` syntax.


```go
//go.mod
module main

go 1.24

// ...
replace github.com/razshare/frizzante => /home/user1/path/to/frizzante
require github.com/razshare/frizzante v1.14.29
// ...
```

This will make it so that changes to your local `frizzante` project are immediately picked up when building your local `frizzante-starter`.

> [!DANGER]
> Another reason for using local packages is to avoid unnecessary mirror caching in https://proxy.golang.org/.
> 
> **Explanation**
> 
> Whenever you request packages for you Go program, these packages are by default fetched using the GO module mirror.
> 
> The Go module mirror automatically caches package requests for 30 minutes, regardless if the requested package or version of the package actually exists or not.
> 
> This means that if you modify, for example
> 
> ```go
> require github.com/razshare/frizzante v1.14.29
> ```
> 
> Into
> 
> ```go
> require github.com/razshare/frizzante v1.14.30
> ```
> 
> And `v1.14.30` doesn't actually exist yet, you package request will be rejected with a `404 Not Found`.\
> That response is then cached for 30 minutes.
> 
> This means that subsequent requests for `v1.14.30` will fail until the cached state is evicted, even if in the meantime `v1.14.30` has actually been published.
> 
> Some IDEs will save your `go.mod` file and try update the package version **automatically**!
> 
> Such an example is [GoLand](https://www.jetbrains.com/go).
> 
> Pointing to a local version of the package using `replace` avoids all this.
> 
> Source - https://proxy.golang.org/#faq-new-version

## Create Branch

Create a new branch and give it a name that describes your changes.

```sh
git checkout -b feature/some-feature
```

## Coding Standards

Submitted code must follow a few rules.

### Functions and Types

Each package must contain a `functions.go` file and/or a `types.go` file.
Functions must be placed in `functions.go` and types must be placed in `types.go`.

```sh
- package1
    - functions.go
    - types.go
```

### Package Nesting

Packages must not be nested deeper than 1 level.

```sh
- package1
    - functions.go
    - types.go
- package2
    - functions.go
    - types.go
- package3
    - functions.go
    - types.go
- package4
    - functions.go
    - types.go
```

### Tests Positioning

Tests must be located at the root of the repository and must reflect the name of the package.

```sh
- package1
    - functions.go
    - types.go
- test_package1.go
```

### Testing Servers

Do not start multiple instances of `*servers.Server` when testing server features.

The whole test suite makes use of one single server instance, which is initialized in [init_test.go](https://github.com/razshare/frizzante/blob/main/init_test.go).

Whenever you need to add more server related tests, make use of the same server instance.

This helps keeping the test suite runtime low.

> [!TIP]
> If you really need to start a new server instance, make sure to use a different port number.

### Encapsulation

Structures and packages **must not** declare private members.

All package functions, variables and structure properties **must always be public**.

```go
//package1/types.go
type MyStruct struct {
  Property1 string
  Property2 int
  Property3 bool
  Property4 any
}
```

```go
//package1/functions.go
func (str *MyStruct) MyFunction1() {}
func (str *MyStruct) MyFunction2() {}
```

## Pull Requests

When you're done with your changes you can submit a pull request in order to implement them into `frizzante` (or `frizzante-starter`).

## Tests Triggers

Tests will automatically run through GitHub actions when pushing into main or when pull requests are opened into main.

That being said, if you don't want to wait for GitHub actions to make sure tests pass, you can also run them locally using the provided git hooks, see next section.

### 

## Git Hooks

You can apply pre-commit git hooks to you local repository by running

```sh
make hooks
```