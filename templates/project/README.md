# What is this?

This is a todo list application

> [!NOTE]
> #### Prerequisites
> Make sure you have `frizzante`, `air`, `bun` and `build-essential` installed on your machine.
>
> ```sh
> sudo apt install build-essential
> which frizzante || go install github.com/razshare/frizzante@latest
> which air || go install github.com/air-verse/air@latest
> which bun || curl -fsSL https://bun.sh/install | bash
> ```

# Get Started

Start development mode with 

```sh
make dev
```

# Build

Build for production with

```sh
make build
```

This will create a standalone `bin/app` binary file.