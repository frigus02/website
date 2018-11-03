# Generator

> Static site generator for this website.

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
