package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
)

func init() {
	tmp = template.Must(template.New("").Parse(defaultTemplate))
}

var tmp *template.Template

var defaultTemplate = `<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title>
            "Choose Your Own Adventure"
        </title>
    </head>
    <body>
        <h1>{{.Title}}</h1>
        {{range .Paragraphs}}
            <p>{{.}}</p>
        {{end}}
        <ul>
            {{range .Options}}
                <li><a href="/{{.Arc}}">{{.Text}}</a></li>
            {{end}}
        </ul>
    </body>
</html>`

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

type handler struct {
	s Story
}

func NewHandler(s Story) http.Handler {
	return handler{s}
}

func (h handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:]
	if chapter, ok := h.s[ArcID(path)]; ok {
		err := tmp.Execute(rw, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(rw, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(rw, "Chapter not found...", http.StatusNotFound)

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
