package main

import (
	"github.com/Imputes/fdlr/internal/errorHandle"
	"github.com/Imputes/fdlr/internal/resume"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(taskCmd)
}

var taskCmd = &cobra.Command{
	Use:     "task",
	Short:   "show current downloading task",
	Example: `fdlr task`,
	Run: func(cmd *cobra.Command, args []string) {
		errorHandle.ExitWithError(task())
	},
}

func task() error {
	err := resume.TaskPrint()
	return errors.WithStack(err)
}
