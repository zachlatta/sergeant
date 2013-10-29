# sergeant

Sergeant is an opinionated command-line application framework for Go.

## Installation

    $ go get github.com/zachlatta/sergeant/sergeant

## Getting Started

Sergeant is very easy to get started with. A basic application can be as simple
as the following:

**main.go**
```go
package main

import "github.com/zachlatta/sergeant/sergeant"

func main() {
  sergeant.ApplicationName = "helloworld"
  sergeant.ApplicationDescription = "helloworld prints `Hello World!` to the console."

  sergeant.Commands = []*sergeant.Command {
    cmdHello,
  }
  sergeant.Init()
}
```

**hello.go**
```go
package main

import (
  "fmt"
  "github.com/zachlatta/sergeant/sergeant"
)

var cmdHello = &sergeant.Command {
  UsageLine: "hello",
  Short: "say hello to the world",
  Long: `
Hello greets the world with a friendly message.
`,
}

func init() {
  cmdHello.run = runHello
}

func runHello(cmd *sergeant.Command, args []string) {
  fmt.Println("Hello world!")
}
```

Running our application produces the following:

    helloworld prints 'Hello World!' to the console.

    Usage:

      helloworld command [arguments]

    The commands are:

      hello       say hello to the world

    Use "helloworld help [command]" for more information about a command.

For more usage examples, check out the `examples` folder.
