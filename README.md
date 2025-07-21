# Frizzante

<a href="https://github.com/razshare/frizzante/releases"><img src="https://img.shields.io/github/release/razshare/frizzante" alt="Latest Release"></a>
<a href="https://github.com/razshare/frizzante/actions"><img src="https://github.com/razshare/frizzante/actions/workflows/tests.yml/badge.svg?branch=main" alt="Build Status"></a>
<a href="https://discord.gg/qgetCNUJ"><img src="https://dcbadge.limes.pink/api/server/https://discord.gg/qgetCNUJ?style=flat" alt="Discord Community"></a>

Frizzante is an opinionated web server framework written in [Go](https://go.dev/) that uses [Svelte](https://svelte.dev/docs/svelte/overview) to render web pages.

# Prerequisites

Install `frizzante`, `air`, `bun`, `build-essential`.

```sh
sudo apt install build-essential
go install github.com/razshare/frizzante@latest
go install github.com/air-verse/air@latest
which bun || curl -fsSL https://bun.sh/install | bash
```

>[!TIP]
>Remember to add Go binaries to your path.
>
> ```sh
> export GOPATH=$HOME/go
> export PATH=$PATH:$GOPATH/bin
> ```

# Get started

Create a project with
```sh
frizzante -c MyProject
```

Install dependencies with

```sh
make install
```

Start development with

```sh
make dev
```

Build with

```sh
make build
```

This will create a `.gen/bin/app` standalone executable.


> [!NOTE]
> The final executable uses [V8](https://v8.dev/) bindings to run
> JavaScript code on the server in order to render svelte components.
>
> For that reason, the first run may take some time, just be patient.
>
> Subsequent runs will take considerably less time.

> [!NOTE]
> Frizzante is aimed mainly at linux distributions.\
> Feel free to contribute improvements for other platforms.
