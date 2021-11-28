package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/Imputes/fdlr/internal/errorHandle"
	"github.com/Imputes/fdlr/internal/executioner"
	"github.com/Imputes/fdlr/internal/tool"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var conc int

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().IntVarP(&conc, "goroutines count", "c", runtime.NumCPU(), "default is your CPU threads count")
}

var downloadCmd = &cobra.Command{
	Use:     "download",
	Short:   "downloads a file from URL or file name",
	Example: `fdlr download [-c=goroutines_count] URL`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		errorHandle.ExitWithError(download(args))
	},
}

func download(args []string) error {
	folder, err := tool.GetFolderFrom(args[0])
	if err != nil {
		return errors.WithStack(err)
	}
	if tool.IsFolderExisted(folder) {
		fmt.Printf("Task already exist, remove it first \n")
		folder, err = tool.GetFolderFrom(args[0])
		if err != nil {
			return errors.WithStack(err)
		}
		if err := os.RemoveAll(folder); err != nil {
			return errors.WithStack(err)
		}
	}

	return executioner.Do(args[0], nil, conc)
}
