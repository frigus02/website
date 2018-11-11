package build

import (
	"html/template"
	"strings"
	"time"

	"github.com/frigus02/website/generator/data"
)

type layoutContext struct {
	ID          string
	Title       string
	Content     template.HTML
	ParentID    string
	ParentTitle string
	StaticFiles map[string]string
}

type pageContext struct {
	Posts       []*data.Post
	Projects    []*data.Project
	StaticFiles map[string]string
}

type dataPageContext struct {
	Item        interface{}
	StaticFiles map[string]string
}

func parseTemplate(content string) (tmpl *template.Template, err error) {
	funcMap := template.FuncMap{
		"hasSuffix": strings.HasSuffix,
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02 15:04")
		},
	}

	tmpl, err = template.New("").Funcs(funcMap).Parse(content)
	return
}
