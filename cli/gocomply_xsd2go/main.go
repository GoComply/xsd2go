package main

import (
	"github.com/gocomply/xsd2go/cli/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
