package data

import "time"

const postsDir = "posts"

// Post is the data model for one post.
type Post struct {
	ID       string
	Title    string    `yaml:"title"`
	Summary  string    `yaml:"summary"`
	Datetime time.Time `yaml:"datetime"`
	Tags     []string  `yaml:"tags"`
	Content  string
}

func (p *Post) setID(id string) {
	p.ID = id
}

func (p *Post) setContent(content string) {
	p.Content = content
}

// ReadPosts reads all posts from the data directory.
func ReadPosts() (*[]Post, error) {
	posts := []Post{}
	err := walkDataDir(postsDir, func(parentDir, itemDir string) error {
		post := Post{}
		err := readDataItem(parentDir, itemDir, &post)
		if err != nil {
			return err
		}

		posts = append(posts, post)
		return nil
	})

	return &posts, err
}
