package main

import (
	"os"

	"github.com/spf13/cobra"
)

// when cli called without any child commands
var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "file downloader written in Go",
}
