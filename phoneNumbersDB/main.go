package main

import (
	"bytes"
	"fmt"
)

func main() {
	fmt.Println("vim-go")
}

func normalize(number string) string {
	var res bytes.Buffer
	for _, ch := range number {
		if ch >= '0' && ch <= '9' {
			res.WriteRune(ch)
		}
	}
	return res.String()
}
