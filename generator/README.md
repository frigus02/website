# Generator

> Static site generator for this website.

# Folder structure

```
(project)
|- data
|  |- posts
|  |  `- 00-my-post-1
|  |     `- index.md          - Metadata and content of the post
|  `- projects
|     `- 20150102-my-project-1
|        |- image1.png        - An image, copied to /static/data with content hash added to filename
|        `- index.md          - Metadata and content of the project
|- pages
|  |- posts
|  |  |- _details.html        - Layout for every post in the data directory
|  |  `- index.html           - Index page of /posts
|  |- projects
|  |  |- _details.html        - Layout for every project in the data directory
|  |  `- index.html           - Index page of /projects
|  |- _layout.html            - Main layout
|  |- 404.html                - Other pages
|  `- index.html              - Index page of the website
|- static                     - Copied to output with content hash added to filename
`- favicon.ico                - A top level file, copied to output
```

# Data model

The layout will be rendered with this model:

```json
{
    "ID": "impress",
    "Title": "Impress",
    "Content": "<section>Address: ...</section>",
    "ParentID": "index",
    "ParentTitle": "About",
    "StaticFiles": {
        "images/me.jpg": "static/images/me-<hash>.jpg",
        "styles/abc.css": "static/styles/abc-<hash>.css"
    }
}
```

Normal pages will be rendered with this model:

```json
{
    "Posts": [
        {
            "ID": "post-1",
            "Order": 39,
            "Metadata": {}
        }
    ],
    "Projects": [
        {
            "ID": "project-1",
            "Order": 20150201,
            "Metadata": {}
        }
    ],
    "StaticFiles": {
        "images/me.jpg": "static/images/me-<hash>.jpg",
        "styles/abc.css": "static/styles/abc-<hash>.css"
    }
}
```

Data pages (pages showing the details of one data entry) have the metadata of their respective data type in the property "Metadata" of the following model:

```json
{
    "ID": "post-1",
	"Order": 39,
	"Metadata": {},
	"Content": "<p>Hello</p>...",
    "StaticFiles": {
        "images/me.jpg": "static/images/me-<hash>.jpg",
        "styles/abc.css": "static/styles/abc-<hash>.css"
    }
}
```
