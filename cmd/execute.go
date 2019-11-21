package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	inputFile string
	chunkSize uint32
	debug     bool
	log       = logrus.New()
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

	rootCmd.PersistentFlags().BoolVarP(
		&debug,
		"debug", "d",
		false,
		"Enable debug mode",
	)

	// Add commands
	rootCmd.AddCommand(sendCmd)
}

var rootCmd = &cobra.Command{
	Use: "splitter",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			log.SetLevel(logrus.DebugLevel)
		}
	},
}

// Execute executes root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Printf("failed to execute root command: %v", err)
		os.Exit(1)
	}
}
