# What is this?

This is a todo list application

# Get Started

Configure project.

    ```sh
    frizzante configure
    ```

Migrate development.

    ```sh
    frizzante migrate
    ```

Start development.

    ```sh
    frizzante dev
    ```

Build production.

    ```sh
    frizzante build
    ```
    
    This will create two executables, `.gen/bin/migrate` and `.gen/bin/serve`.

Migrate production.

    ```sh
    .gen/bin/migrate
    ```

Serve production.

    ```sh
    .gen/bin/serve
    ```
> [!NOTE]
> `.gen/bin/serve` is a standalone binary, it contains all html, css and javascript bundles.
