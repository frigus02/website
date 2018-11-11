package data

import (
	"html/template"
	"time"
)

const postsDir = "posts"

// Post is the data model for one post.
type Post struct {
	ID       string
	Order    int
	Title    string    `yaml:"title"`
	Summary  string    `yaml:"summary"`
	Datetime time.Time `yaml:"datetime"`
	Tags     []string  `yaml:"tags"`
	Content  template.HTML
}

func (p *Post) setID(id string) {
	p.ID = id
}

func (p *Post) setOrder(order int) {
	p.Order = order
}

func (p *Post) setContent(content template.HTML) {
	p.Content = content
}
