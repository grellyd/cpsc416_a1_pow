// TODO: Documentation Here
package main

import (
	"a1/cli"
	"fmt"
	"os"
)

func main() {
	exitCode, err := cli.Run(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	os.Exit(exitCode)
}
