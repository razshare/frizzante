# Frizzante

<img src="https://raw.githubusercontent.com/razshare/frizzante/refs/heads/main/assets/frizz-octo-header.webp" width="308" />

<a href="https://github.com/razshare/frizzante/releases"><img src="https://img.shields.io/github/release/razshare/frizzante" alt="Latest Release"></a>
<a href="https://github.com/razshare/frizzante/actions"><img src="https://github.com/razshare/frizzante/actions/workflows/tests.yaml/badge.svg?branch=main" alt="Tests Status"></a>
<a href="https://discord.gg/y7tTeR7yPH"><img src="https://dcbadge.limes.pink/api/server/https://discord.gg/y7tTeR7yPH?style=flat" alt="Discord Community"></a>

Frizzante is an opinionated web server framework written in [Go](https://go.dev/) that uses [Svelte](https://svelte.dev/docs/svelte/overview) to render web pages.

# Prerequisites

Install frizzante.

```sh
go run github.com/razshare/frizzante/setup/install@latest
```

> [!NOTE]
> Remember to add Go binaries to your path.
> 
> ```sh
> export GOPATH=$HOME/go
> export PATH=$PATH:$GOPATH/bin
> ```

# Get Started

Create project.
```sh
frizzante -c MyProject
```

Configure project.

```sh
frizzante --configure
```

Start development.

```sh
frizzante --dev
```

Build.

```sh
frizzante --build
```

This will create a `.gen/bin/app` standalone executable.

# Thanks

Thanks to [cmjoseph07](https://github.com/cmjoseph07) for the octo mascot!
