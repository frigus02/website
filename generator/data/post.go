package data

import (
	"time"
)

// PostsDir is the directly name under "data/", which contains posts.
const PostsDir = "posts"

// PostMetadata contains the metadata for the YAML header stored in post
// markdown files.
type PostMetadata struct {
	Title    string    `yaml:"title"`
	Summary  string    `yaml:"summary"`
	Datetime time.Time `yaml:"datetime"`
	Tags     []string  `yaml:"tags"`
}
