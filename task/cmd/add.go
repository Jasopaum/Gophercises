package cmd

import (
    "fmt"
    "strings"
    "github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
    Use:   "add",
    Short: "Add a task to my TODO list.",
    Run: func(cmd *cobra.Command, args []string) {
        task := strings.Join(args, " ")
        fmt.Printf("Added %s to list.\n", task)
    },
}

func init() {
    RootCmd.AddCommand(addCmd)
}
