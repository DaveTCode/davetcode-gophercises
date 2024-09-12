package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Story map[string]StoryArc

type StoryArc struct {
	Title      string        `json:"title"`
	Paragraphs []string      `json:"story"`
	Options    []StoryOption `json:"options"`
}

type StoryOption struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func ParseStoryJson(storyFilePath string) (Story, error) {
	storyBytes, err := os.ReadFile(storyFilePath)

	if err != nil {
		log.Fatal("Can't open file "+storyFilePath, err)
	}

	var story Story
	if err := json.Unmarshal(storyBytes, &story); err != nil {
		return nil, err
	}

	return story, nil
}

func CreateHttpStoryHandler(story Story, mux http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		arc, exists := story[r.URL.Path[1:]]

		if !exists {
			http.Redirect(w, r, "/intro", http.StatusFound)
		} else {
			err := storyArcToHtml(arc, w)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}
		}
	})
}

func storyArcToHtml(storyArc StoryArc, w http.ResponseWriter) error {
	const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		{{range .Paragraphs}}<p>{{ . }}</p>{{else}}<p><strong>Blank arc</strong></p>{{end}}
		<ul>{{range .Options}}<li><a href="/{{ .Arc }}">{{ .Text }}</a></li>{{else}}<li><strong>End of arc</strong></li>{{end}}</ul>
	</body>
</html>`

	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		return err
	}

	t.Execute(w, storyArc)

	return nil
}

func main() {
	var storyFileName string
	flag.StringVar(&storyFileName, "story", "gopher.json", "The json story file")

	story, err := ParseStoryJson(storyFileName)
	if err != nil {
		panic(err)
	}

	mux := defaultMux()
	storyHandler := CreateHttpStoryHandler(story, mux)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", storyHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", redirectToIntro)
	return mux
}

func redirectToIntro(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/intro", http.StatusFound)
}
