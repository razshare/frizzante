<img alt="frizzante logo" src="https://raw.githubusercontent.com/razshare/frizzante/refs/heads/main/assets/logo.png" width="308" />

<a href="https://github.com/razshare/frizzante/v2/releases"><img src="https://img.shields.io/github/release/razshare/frizzante" alt="Latest Release"></a>
<a href="https://github.com/razshare/frizzante/v2/actions"><img src="https://github.com/razshare/frizzante/actions/workflows/tests.yaml/badge.svg?branch=main" alt="Tests Status"></a>
<a href="https://discord.gg/y7tTeR7yPH"><img src="https://img.shields.io/badge/Discord-Join%20Server-7289da?logo=discord&logoColor=white&style=flat" alt="Discord Community"></a>

Frizzante is an opinionated web server framework written in [Go](https://go.dev/) that uses [Svelte](https://svelte.dev/docs/svelte/overview) to render web pages.

# Prerequisites

Install frizzante.

```sh
go install github.com/razshare/frizzante/v2@latest
```

> [!TIP]
> Remember to add Go binaries to your path.
> 
> ```sh
> export GOPATH=$HOME/go
> export PATH=$PATH:$GOPATH/bin
> ```

# Get Started

1. Create project.
    
```sh
frizzante create my_project
```

2. Configure project.
    
```sh
frizzante configure
```

3. Migrate development.
    
```sh
frizzante migrate
```

4. Start development.
    
```sh
frizzante dev
```

5. Build production.
    
```sh
frizzante build
```

This will create two executables, `.gen/bin/migrate` and `.gen/bin/serve`.

6. Migrate production.
    
```sh
.gen/bin/migrate
```

7. Serve production.
    
```sh
.gen/bin/serve
```

> [!NOTE]
> `.gen/bin/serve` is a standalone binary, it contains all html, css and javascript bundles.
