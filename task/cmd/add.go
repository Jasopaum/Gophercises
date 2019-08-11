package cmd

import (
    "fmt"
    "strings"

    "gophercises/task/db"

    "github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
    Use:   "add",
    Short: "Add a task to my TODO list.",
    Run: func(cmd *cobra.Command, args []string) {
        task := strings.Join(args, " ")
        err := db.CreateTask(task)
        if err != nil {
            fmt.Printf("Error: ", err)
        }
    },
}

func init() {
    RootCmd.AddCommand(addCmd)
}
