package main

import (
	"os"

	"github.com/alecthomas/kingpin/v2"
)

var app = kingpin.New("cli", "A command-line interface")

var commandSomething = app.Command("something", "Do something")
var somethingX = commandSomething.Flag("x", "X").Short('x').String()
var somethingY = commandSomething.Flag("y", "Y").Short('y').String()

var commandAdd = app.Command("add", "Add User")
var commandAddArgsUser = commandAdd.Arg("user", "username").Required().String()
var commandAddFlagOverride = commandAdd.Flag("override", "Override").Short('o').Bool()

var commandDelete = app.Command("delete", "Delete User")
var commandDeleteForce = commandDelete.Flag("force", "Force").Short('f').Bool()
var commandDeleteArgs = commandDelete.Arg("user", "username").Required().String()

func main() {

	commandAdd.Action(func(c *kingpin.ParseContext) error {
		println("add", *commandAddArgsUser, *commandAddFlagOverride)
		return nil
	})

	commandDelete.Action(func(c *kingpin.ParseContext) error {
		println("delete", *commandDeleteArgs, *commandDeleteForce)
		return nil
	})

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
