package data

import "html/template"

const projectsDir = "projects"

// Project is the data model for one project.
type Project struct {
	ID               string
	Title            string          `yaml:"title"`
	ShortDescription string          `yaml:"short_description"`
	Images           []ProjectImage  `yaml:"images"`
	Sources          []ProjectSource `yaml:"sources"`
	Tags             []string        `yaml:"tags"`
	Content          template.HTML
}

func (p *Project) setID(id string) {
	p.ID = id
}

func (p *Project) setContent(content template.HTML) {
	p.Content = content
}

// ProjectImage is the data model for images referenced in a Project.
type ProjectImage struct {
	Name string `yaml:"name"`
	Alt  string `yaml:"alt"`
}

// ProjectSource is the data model for sources referenced in a Project.
type ProjectSource struct {
	Type ProjectSourceType `yaml:"type"`
	URL  string            `yaml:"url"`
}

// ProjectSourceType is the type for project sources.
type ProjectSourceType string

// Project source types are available as constants.
const (
	GitProjectSource            ProjectSourceType = "git"
	GooglePlayProjectSource     ProjectSourceType = "googleplay"
	AMOProjectSource            ProjectSourceType = "amo"
	ChromeWebStoreProjectSource ProjectSourceType = "chromewebstore"
	TryProjectSource            ProjectSourceType = "try"
)
