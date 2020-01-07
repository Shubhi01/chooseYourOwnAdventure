package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	//port := flag.Int("port", 3030, "Port to access server")
	storyName := flag.String("storyName", "gopher.json", "Name of file to read story from")
	flag.Parse()

	// Read story content from the file storyName
	content, err := ioutil.ReadFile(*storyName)
	if err != nil {
		fmt.Println("ERROR: Failed to read from story file")
		os.Exit(1)
	}

	story := NewStoryFromJSON(content)
	c := NewCLIHandler(story)
	c.StartStory()
	// h := NewHandler(story)
	// fmt.Println("Starting server at port %d", *port)
	// log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
	//fmt.Println(c)
}
