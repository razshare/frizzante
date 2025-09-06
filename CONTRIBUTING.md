This document describes the full process of setting up a fully working local development environment and submitting your first contribution.

## Clone Repository

Clone the `frizzante` repository from GitHub.

```go
git clone https://github.com/razshare/frizzante
```

> [!TIP]
> If you don't have direct access to it you will need to fork your own `frizzante` repository.
> 
> Then when you're done with your changes you will need to submit a pull request.

## Install Frizzante

Install frizzante.

```sh
go run github.com/razshare/frizzante/install@v1.26.6
```

> [!TIP]
> Remember to add Go binaries to your path.
> 
> ```sh
> export GOPATH=$HOME/go
> export PATH=$PATH:$GOPATH/bin
> ```

## Create Branch

Create a new branch and give it a name that describes your changes.

```sh
git checkout -b feature/some-feature
```

## Coding Standards

Submitted code must follow a few rules.


### Export Everything

All package functions, variables and structure fields **must always be exported**.

```go
//package1/types.go
type MyStruct struct {
  Field1 string
  Field2 int
  Field3 bool
  Field4 any
}
func MyFunction(){}
const Planet = "Earth"
var Name = "World"
```

### Data & Logic

Data and logic must be separated, don't **define receiver functions** unless you're forced
to by some third party package.

This type of code should be avoided

```go
type MyStruct struct {
  field1 string
}

func (str *MyStruct) MyFunction(){
  str.field1 = "Hello!"
}
```

and be instead converted into

```go
type MyStruct struct {
  Field1 string
}

func MyFunction(str *MyStruct){
  str.Field1 = "Hello!"
}
```

## Pull Requests

When you're done with your changes you can submit a pull request in order to implement them into `frizzante`.

## Testing

### Automation

Tests will automatically run through GitHub actions when pushing into main or when pull requests are opened into main.