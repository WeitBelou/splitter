package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	inputFile string
	chunkSize uint32
)

func init() {
	// Parse flags
	rootCmd.PersistentFlags().StringVarP(
		&inputFile,
		"input-file", "i",
		"-",
		"Input file use '-' for STDIN",
	)

	rootCmd.PersistentFlags().Uint32VarP(
		&chunkSize,
		"chunk-size", "s",
		100,
		"Size of chunks on which file will be split",
	)

	// Add commands
	rootCmd.AddCommand(sendCmd)
}

var rootCmd = &cobra.Command{
	Use: "splitter",
}

// Execute executes root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Printf("failed to execute root command: %v", err)
		os.Exit(1)
	}
}

func getInputFile(name string) (*os.File, error) {
	if name == "-" {
		return os.Stdin, nil
	}
	return os.Open(name)
}
