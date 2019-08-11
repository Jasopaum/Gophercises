package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
    Use:   "list",
    Short: "Display my TODO list.",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Called list task.")
    },
}

func init() {
    RootCmd.AddCommand(listCmd)
}
