/*
Copyright © 2023 jaronnie jaron@jaronnie.com
*/
package main

import (
	"github.com/jaronnie/gvm/cmd"
)

var version = "v1.8.0"

func main() {
	cmd.Version = version

	cmd.Execute()
}
