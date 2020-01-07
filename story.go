package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

// func init() {
// 	tmp = template.Must(template.New("").Parse(defaultTemplate))
// }

func initCLITemplate() {
	fmt.Println("init")
	tmp = template.Must(template.New("").Parse(defaultCLITemplate))
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

var defaultCLITemplate = `{{.Title}}

	{{range .Paragraphs}}
		{{.}}
	{{end}}
	{{range $index, $opt := .Options}}
		Press {{$index}} to venture into {{$opt.Text}}
	{{end}}
`

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
	t *template.Template
}

type HandlerOption func(h *handler)

// If we want to provide options to handler like a custom template,
// custom parsing function, the options are:
// 1. Pass individual options as arguments to NewHander func. This could
// 	become cumbersome
// 2. Create a struct which encapsulates all handler options and pass it as
// 	an argument to NewHandler. It's difficult to define complex dependencies
// 	between different options using this method
// 3. Functional arguments return a function which takes pointer to handler
// as argument and updates it with the passed option. (using closures)

// WithTemplate is an example of using functional options
func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tmp}
	for _, opt := range opts {
		opt(&h)
	}
	return h
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

type CLIHandler struct {
	s Story
}

func NewCLIHandler(s Story) CLIHandler {
	return CLIHandler{s}
}

func (c CLIHandler) StartStory() {
	initCLITemplate()
	//tmp.Execute(os.Stdout, c.s["intro"])

	chapter := "intro"
	for {
		tmp.Execute(os.Stdout, c.s[ArcID(chapter)])
		var input int
		fmt.Scanf("%d", &input)
		chapter = string(c.s[ArcID(chapter)].Options[input-1].Arc)
	}
}
