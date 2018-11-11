package build

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/frigus02/website/generator/data"
	"github.com/frigus02/website/generator/fs"
	"golang.org/x/tools/blog/atom"
)

type feedItem struct {
	dependencies []item
}

func newFeedItem() *feedItem {
	return &feedItem{}
}

func (i *feedItem) addItem(item item) {
	switch item := item.(type) {
	case *dataItem:
		switch item.metadata.(type) {
		case *data.PostMetadata:
			i.dependencies = addItemIfNotExists(i.dependencies, item)
		}
	}
}

func (i *feedItem) isDependentOn(item item) bool {
	return containsItem(i.dependencies, item)
}

func (i *feedItem) update(file *fs.File) error {
	return nil
}

func (i *feedItem) render(ctx *renderContext) error {
	var posts []*dataItem
	for _, item := range i.dependencies {
		if post, ok := item.(*dataItem); ok {
			posts = append(posts, post)
		}
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].order > posts[j].order
	})

	baseURL := "https://kuehle.me"
	author := &atom.Person{
		Name:  "Jan Kuehle",
		URI:   baseURL,
		Email: "jkuehle90@gmail.com",
	}
	feed := &atom.Feed{
		Title: "Jan Kuehle - Blog",
		ID:    fmt.Sprintf("%s/posts", baseURL),
		Link: []atom.Link{
			atom.Link{
				Rel:  "alternate",
				Href: fmt.Sprintf("%s/posts", baseURL),
			},
			atom.Link{
				Rel:  "self",
				Href: fmt.Sprintf("%s/feeds/posts", baseURL),
			},
		},
		Updated: atom.Time(time.Now()),
		Author:  author,
	}

	for _, post := range posts {
		if metadata, ok := post.metadata.(*data.PostMetadata); ok {
			feed.Entry = append(feed.Entry, &atom.Entry{
				Title: metadata.Title,
				ID:    fmt.Sprintf("%s/posts/%s", baseURL, post.id),
				Link: []atom.Link{
					atom.Link{
						Href: fmt.Sprintf("%s/posts/%s", baseURL, post.id),
					},
				},
				Published: atom.Time(metadata.Datetime),
				Updated:   atom.Time(metadata.Datetime),
				Author:    author,
				Summary: &atom.Text{
					Type: "html",
					Body: metadata.Summary,
				},
				Content: &atom.Text{
					Type: "html",
					Body: string(post.content),
				},
			})
		}
	}

	data, err := xml.Marshal(feed)
	if err != nil {
		return fmt.Errorf("error marshaling feed to XML: %v", err)
	}

	filename := filepath.Join(ctx.out, "feeds/posts")
	dir := filepath.Dir(filename)

	err = os.MkdirAll(dir, 0644)
	if err != nil {
		return fmt.Errorf("error creating destination folder %s for feeds: %v", dir, err)
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("error creating feed: %v", err)
	}

	return nil
}
