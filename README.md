# Get started

Clone the basic starter template.
```sh
git clone https://github.com/razshare/frizzante-starter && \
rm frizzante-starter/.git -fr && \
cd frizzante-starter
```

> [!NOTE]
> Make sure [Bun](https://bun.sh/) and [build-essential](https://askubuntu.com/questions/398489/how-to-install-build-essential) are installed.

Start development with

```sh
make dev
```

> [!TIP]
> The default makefile uses [Bun](https://bun.sh/) to install 
> dependencies while in development mode.\
> \
> You can swap Bun out with whatever runtime you want to use to 
> develop by replacing all `bun` and `bunx` references within 
> the makefile with the equivalent of whatever runtime you'd like to use.\
> \
> Note that this will have no impact on your application's performance, 
> Bun is only used in development mode
> in order to be able to run the Vite server.

Build with

```sh
make build
```

This will create a `bin/app` standalone executable.

> The final executable uses [V8](https://v8.dev/) bindings to run JavaScript code on the server in order to render svelte components.\
> \
> For that reason, the first build or launch may take some time, just be patient.\
> Once that is done, subsequent builds or launches will take considerably less time.