# Get started

Frizzante is an opinionated web server framework that uses Svelte to render web pages.

Make sure you have [Go](https://go.dev/doc/install),
[build-essential](https://askubuntu.com/questions/398489/how-to-install-build-essential) and
[Bun](https://bun.sh/) installed.

> [!NOTE]
> Frizzante is aimed mainly at linux distributions.

Create a new project by cloning the started template.

```sh
git clone https://github.com/razshare/frizzante-starter && \
rm frizzante-starter/.git -fr
```

Configure the project

```sh
cd frizzante-starter && \
make configure
```

Start the server in development mode with

```sh
make dev
```

Build the project it into a standalone executable for production with

```sh
make build
```

> [!NOTE]
> The final output is a standalone executable.