package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	ver = "dev"
)

var rootCmd = &cobra.Command{
	Use:     "gowasm",
	Short:   "gowasm",
	Version: ver,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
