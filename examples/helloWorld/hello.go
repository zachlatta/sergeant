package main

import (
  "fmt"
  "github.com/zachlatta/sergeant/sergeant"
)

var cmdHello = &sergeant.Command{
  UsageLine: "hello",
  Short: "say hello to the world",
  Long: `
Hello greets the world with a friendly message.

Example usage:

  helloWorld hello

`,
}

func init() {
  cmdHello.Run = runHello
}

func runHello(cmd *sergeant.Command, args []string) {
  fmt.Println("Hello world!")
}
