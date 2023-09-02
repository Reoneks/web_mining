package main

import (
	"dyploma/cmd"

	"go.uber.org/fx"
)

func main() {
	fx.New(cmd.Exec()).Run()
}
