package main

import (
	"encoding/json"
	"fmt"
)

// ArcID is unique key for a story arc
type ArcID string

// Story consists of a number of story arcs
type Story map[ArcID]Arc

type Arc struct {
	Title      string
	Paragraphs []string `json:"story"`
	Options    []Options
}

type Options struct {
	Text string
	Arc  ArcID
}

// NewStoryFromJSON returns a new story from json bytes
func NewStoryFromJSON(content []byte) Story {
	story := Story{}
	err := json.Unmarshal(content, &story)
	if err != nil {
		fmt.Println("ERROR: JSON unmarshal failed with error ", err)
	}

	return story
}
