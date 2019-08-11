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
        tasks, err := db.ListTasks()
        if err != nil {
            fmt.Println("Error: ", err)
        }
        for i, t := range tasks {
            fmt.Printf("%d: %s\n", i+1, t.Value)
        }
    },
}

func init() {
    RootCmd.AddCommand(listCmd)
}
