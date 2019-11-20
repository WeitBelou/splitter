package cmd

import (
	"os"
)

func getInputFile(name string) (*os.File, error) {
	if name == "-" {
		return os.Stdin, nil
	}
	return os.Open(name)
}
