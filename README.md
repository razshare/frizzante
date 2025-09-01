# Frizzante

<img src="https://raw.githubusercontent.com/razshare/frizzante/refs/heads/main/assets/frizz-octo-header.webp" width="308" />

<a href="https://github.com/razshare/frizzante/releases"><img src="https://img.shields.io/github/release/razshare/frizzante" alt="Latest Release"></a>
<a href="https://discord.gg/y7tTeR7yPH"><img src="https://dcbadge.limes.pink/api/server/https://discord.gg/y7tTeR7yPH?style=flat" alt="Discord Community"></a>

Frizzante is an opinionated web server framework written in [Go](https://go.dev/) that uses [Svelte](https://svelte.dev/docs/svelte/overview) to render web pages.

# Prerequisites

Download the [latest binary](https://github.com/razshare/frizzante/releases).

> [!TIP]
> Or build it yourself.
> ```sh
> git clone https://github.com/razshare/frizzante --depth=1
> cd frizzante && make zip && go build -o frizzante
> ```

> [!TIP]
> Remember to add the binary to your path.
> ```sh
> export PATH="/path/to/frizzante:$PATH"
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
