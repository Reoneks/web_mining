package main

import (
	"test/cmd"

	"go.uber.org/fx"
)

func main() {
	fx.New(cmd.Exec()).Run()
}
