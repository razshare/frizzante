# What is this?

This is a todo list application

> [!NOTE]
> #### Prerequisites
> Install `build-essential` and `frizzante`.
>
> ```sh
> sudo apt install build-essential
> go install github.com/razshare/frizzante@latest
> ```

# Get Started

Configure project

```sh
make configure
```

Start development mode with 

```sh
make dev
```

# Build

Build for production with

```sh
make build
```

This will create a standalone `.gen/bin/app` binary file.
