package cmd

import (
    "fmt"
    "strconv"

    "gophercises/task/db"

    "github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
    Use:   "do",
    Short: "Set a task of my TODO list to done.",
    Run: func(cmd *cobra.Command, args []string) {
        var idsToDel []int
        tasks, err := db.ListTasks()

        for _, arg := range args {
            id, err := strconv.Atoi(arg)
            if err != nil {
                fmt.Println("Failed to parse: ", arg)
            } else {
                if 0 < id && id <= len(tasks) {
                    idsToDel = append(idsToDel, id-1)
                } else {
                    fmt.Printf("%d not in list.\n", id)
                }
            }
        }

        if err != nil {
            fmt.Println("Error: ", err)
        }
        for _, idToDel := range idsToDel {
            keyToDel := tasks[idToDel].Key
            db.DeleteTask(keyToDel)
        }
    },
}

func init() {
    RootCmd.AddCommand(doCmd)
}
