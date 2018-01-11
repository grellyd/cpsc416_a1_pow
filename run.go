// TODO: Documentation Here
package main

import (
	"os"
	"fmt"
	"a1/cli"
)
func main() {
	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	os.Exit(exitCode)
}
