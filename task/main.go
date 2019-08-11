package main

import (
    "gophercises/task/cmd"
    "gophercises/task/db"
)

func main () {
    db.Init()
    cmd.RootCmd.Execute()
}

