package main

import (
	"os"

	foundry "github.com/jabahum/foundry"
	"github.com/jabahum/foundry/internal/cli"
)

func main() { os.Exit(cli.App{FS: foundry.Templates}.Run(os.Args[1:])) }
