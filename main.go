package main

import (
	_ "go.uber.org/automaxprocs"

	"mybase/cmd/bootstrap"
	_ "mybase/cmd/cli"
)

func main() {
	bootstrap.Run()
}
