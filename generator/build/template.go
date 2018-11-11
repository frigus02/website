package build

import (
	"html/template"
	"strings"

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
	}

	tmpl, err = template.New("").Funcs(funcMap).Parse(content)
	return
}
