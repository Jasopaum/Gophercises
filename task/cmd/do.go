package cmd

import (
    "fmt"
    "strconv"
    "github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
    Use:   "do",
    Short: "Set a task of my TODO list to done.",
    Run: func(cmd *cobra.Command, args []string) {
        var tasks []int
        for _, arg := range args {
            id, err := strconv.Atoi(arg)
            if err != nil {
                fmt.Println("Failed to parse: ", arg)
            } else {
                tasks = append(tasks, id)
            }
        }
        fmt.Println("Did tasks: ", tasks)
    },
}

func init() {
    RootCmd.AddCommand(doCmd)
}
