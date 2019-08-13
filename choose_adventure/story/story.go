package story

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var tpl *template.Template

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title>Choose Your Own Adventure</title>
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

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func WithParsePathFn(f func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = f
	}
}

func defaultParsePath(r *http.Request) string {
	path := r.URL.Path
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:]
}

func NewHandler(st Story, opts ...HandlerOption) http.Handler {
	h := handler{st, tpl, defaultParsePath}
	for _, o := range opts {
		o(&h)
	}
	return h
}

type handler struct {
	s      Story
	t      *template.Template
	pathFn func(r *http.Request) string
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)
	chapter, ok := h.s[path]
	if ok {
		err := tpl.Execute(w, chapter)
		if err != nil {
			http.Error(w, "Error while executing html template", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Chapter not found", http.StatusNotFound)
	}
}

func ReadStoryJson(jsonPath string) Story {
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		fmt.Println("Error while reading json file: ", err)
	}
	defer jsonFile.Close()

	byteStory, _ := ioutil.ReadAll(jsonFile)
	var story Story
	json.Unmarshal([]byte(byteStory), &story)

	return story
}

type Story map[string]Arc

type Arc struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}
