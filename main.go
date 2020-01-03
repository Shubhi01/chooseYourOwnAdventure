package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	var storyName *string
	storyName = flag.String("storyName", "gopher.json", "Name of file to read story from")
	flag.Parse()

	// Read story content from the file storyName
	content, err := ioutil.ReadFile(*storyName)
	if err != nil {
		fmt.Println("ERROR: Failed to read from story file")
		os.Exit(1)
	}

	story := NewStoryFromJSON(content)
	fmt.Println(story)
}
