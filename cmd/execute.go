package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	inputFile string
	outputDir string
	chunkSize uint32
)

func init() {
	// Parse flags
	rootCmd.PersistentFlags().StringVarP(&inputFile, "input-file", "i", "-",
		"Input file (default is '-' which stands for STDIN)",
	)
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output-dir", "o",
		"out",
		"Path where chunks will be saved (default is 'out'",
	)
	rootCmd.PersistentFlags().Uint32VarP(&chunkSize, "chunk-size", "s", 100,
		"Size of chunks on which file will be split (default is '100')",
	)
}

var rootCmd = &cobra.Command{
	Use:   "splitter",
	Short: "Splitter is an example utility that splits big files to chunks of 'n' lines",
	Long:  "Splitter is an example utility that splits big files to chunks of 'n' lines",
	Run:   split,
}

// Execute executes root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Printf("failed to execute root command: %v", err)
		os.Exit(1)
	}
}
