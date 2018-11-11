package data

// ProjectsDir is the directly name under "data/", which contains projects.
const ProjectsDir = "projects"

// ProjectMetadata contains the metadata for the YAML header stored in project
// markdown files.
type ProjectMetadata struct {
	Title            string          `yaml:"title"`
	ShortDescription string          `yaml:"short_description"`
	Images           []ProjectImage  `yaml:"images"`
	Sources          []ProjectSource `yaml:"sources"`
	Tags             []string        `yaml:"tags"`
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
