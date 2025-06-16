# What is this?
Frizzante is a minimalistic and opinionated web server framework written in [Go](https://go.dev/) that uses [Svelte](https://svelte.dev/docs/svelte/overview) to render pages.

# Prerequisites

Make sure you have `build-essential`, `curl` and `unzip` installed on your machine.

```sh
sudo apt install build-essential unzip curl
```

# Get started

Clone the basic starter template.
```sh
git clone https://github.com/razshare/frizzante-starter && \
rm frizzante-starter/.git -fr && \
cd frizzante-starter
```

Update dependencies with

```sh
make update
```

Start development with

```sh
make dev
```

> [!NOTE]
> The default makefile uses [Bun](https://bun.sh/) to install
> dependencies while in development mode.
> 
> You can swap Bun out with whatever runtime you want to use to develop 
> by replacing all `bun` references within the makefile with 
> the equivalent of whatever runtime you'd like to use.
> 
> This change will have no impact on your application's performance, 
> Bun is only used to compile `.svelte` files into `.js` files 
> in development mode or when building for production.

Build with

```sh
make build
```

This will create a `bin/app` standalone executable.


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