package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// set through ldflags
	version string
	commit  string
	date    string
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "prints meta info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:    ", version)
		fmt.Println("Commit:     ", commit)
		fmt.Println("Date:       ", date)
	},
}
