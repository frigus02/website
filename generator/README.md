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
|     `- 00-my-project-1
|        |- icon_large.png    - Large icon for the project
|        |- icon.png          - Icon for the project
|        |- image1_thumb.png  - Thumbnail for the image "image1"
|        |- image1.png        - An image, referenced from metadata by name "image1" (without extension)
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
`- static                     - Copied to output (css is concatenated and minified)
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
    "Stylesheet": "fb5jh3g5j3.css"
}
```

Normal pages will be rendered with this model:

```json
{
    "Posts": [],
    "Projects": []
}
```

Data pages (pages showing the details of one data entry) have the model of their respective data type.
