package main

import (
  "github.com/zachlatta/sergeant/sergeant"
)

func main() {
  sergeant.ApplicationName = "helloworld"
  sergeant.ApplicationDescription = "Prints 'Hello World!' to the console."

  sergeant.Commands = []*sergeant.Command{
    cmdHello,
  }
  sergeant.Init()
}
