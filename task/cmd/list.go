package cmd

import (
    "fmt"

    "gophercises/task/db"

    "github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
    Use:   "list",
    Short: "Display my TODO list.",
    Run: func(cmd *cobra.Command, args []string) {
        db.ListTasks()
        fmt.Println("Called list task.")
    },
}

func init() {
    RootCmd.AddCommand(listCmd)
}
