package data

import (
	"html/template"
	"time"
)

const postsDir = "posts"

// Post is the data model for one post.
type Post struct {
	ID       string
	Title    string    `yaml:"title"`
	Summary  string    `yaml:"summary"`
	Datetime time.Time `yaml:"datetime"`
	Tags     []string  `yaml:"tags"`
	Content  template.HTML
}

func (p *Post) setID(id string) {
	p.ID = id
}

func (p *Post) setContent(content template.HTML) {
	p.Content = content
}

// GetAllPosts reads all posts from the data directory.
func GetAllPosts() (*[]Post, error) {
	posts := []Post{}
	err := walkDataDir(postsDir, func(path, itemDir string) error {
		post := Post{}
		err := readDataItem(path, itemDir, &post)
		if err != nil {
			return err
		}

		posts = append(posts, post)
		return nil
	})

	return &posts, err
}
