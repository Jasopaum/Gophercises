package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"gophercises/choose_adventure/story"
)

func main() {
	var port = flag.Int("port", 8000, "port to start the web app on")
	var storyJson = flag.String("file", "gopher.json", "json file containing the story")

	st := story.ReadStoryJson(*storyJson)

	h := story.NewHandler(st)

	fmt.Printf("Starting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
