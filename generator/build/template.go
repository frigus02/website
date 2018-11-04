package build

import (
	"html/template"

	"github.com/frigus02/website/generator/data"
)

type layoutContext struct {
	ID          string
	Title       string
	Content     template.HTML
	ParentID    string
	ParentTitle string
	Stylesheet  string
}

type pageContext struct {
	Posts    []*data.Post
	Projects []*data.Project
}
