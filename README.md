# Get started

Create a new project using the starter template.

```bash
git clone https://github.com/razshare/frizzante-starter && cd frizzante-starter && rm .git -fr
```

> [!NOTE]
> Make sure you have [Go](https://go.dev/doc/install) and [Bun](https://bun.sh) installed.\
> If you'd rather use a different runtime than Bun to update your javascript dependencies, 
> see [makefile](https://github.com/razshare/frizzante-starter/blob/main/makefile) and change all occurrences of 
> `bun` and `bunx` with the equivalent of whatever runtime you'd like to use.

Update dependencies

```bash
make update
```

> [!NOTE]
> Make sure you have `build-essential` installed
> ```bash
> sudo apt install build-essential
> ```

Then start the server

```bash
make start
```

or build it

```bash
make build
```

> [!NOTE]
> The `.dist` directory is embedded, which makes the final executable completely portable.
> 
> That being said, you can still create a ".dist" directory near your executable.\
> Whenever the server will try to access a file missing from the embedded file system, the server will fall 
> back to the nearby ".dist" directory instead.

> [!NOTE]
> This project is aimed mainly at linux distributions.\
> Feel free to contribute any fixes for other platforms.

