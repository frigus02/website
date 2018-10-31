/* eslint-env node */

"use strict";

const fs = require("fs-extra");
const { Feed } = require("feed");
const marked = require("marked");

async function createIndex(folder, metadataKeys, reverse) {
    await fs.remove(`src/api/${folder}`);
    await fs.mkdirs(`src/api/${folder}`);

    const files = await fs.readdir(`data/${folder}`);
    let index = [];

    for (const fileName of files) {
        const { metadata } = await getMetadata(`data/${folder}`, fileName);
        index.push(metadata);

        await fs.copy(
            `data/${folder}/${fileName}`,
            `src/api/${folder}/${metadata.id}.md`
        );
    }

    index = index.map(metadata => filterMetadata(metadata, metadataKeys));
    if (reverse) {
        index.reverse();
    }

    const indexJson = JSON.stringify(index);
    await fs.writeFile(`src/api/${folder}/index.json`, indexJson, "utf-8");
}

async function createFeed() {
    await fs.remove("src/feeds");
    await fs.mkdirs("src/feeds");

    const baseUrl = "https://kuehle.me";
    const author = {
        name: "Jan Kuehle",
        email: "jkuehle90@gmail.com",
        link: baseUrl
    };

    const feed = new Feed({
        title: "Jan Kuehle - Blog",
        description:
            "Everything I stumble across while coding on my projects. Will propably be something about web or android development.",
        id: `${baseUrl}/posts`,
        link: `${baseUrl}/posts`,
        feedLinks: {
            atom: `${baseUrl}/feeds/posts`
        },
        author: author
    });

    const files = await fs.readdir("data/posts");
    for (const fileName of files) {
        const { metadata, content } = await getMetadata("data/posts", fileName);

        feed.addItem({
            title: metadata.title,
            id: `${baseUrl}/posts/${metadata.id}`,
            link: `${baseUrl}/posts/${metadata.id}`,
            description: metadata.summary,
            content: marked(content),
            author: [author],
            published: new Date(metadata.datetime),
            date: new Date(metadata.datetime)
        });
    }

    const atom = feed.atom1();
    await fs.writeFile("src/feeds/posts", atom, "utf-8");
}

async function getMetadata(folder, fileName) {
    const text = await fs.readFile(`${folder}/${fileName}`, "utf-8");

    const metadataStart = text.indexOf("```json") + 7;
    const metadataEnd = text.indexOf("```", metadataStart);
    const metadataJson = text.substring(metadataStart, metadataEnd);
    const metadata = JSON.parse(metadataJson);
    metadata.id = fileName.split("-")[1].split(".")[0];

    const content = text.substring(metadataEnd + 3);

    return {
        metadata,
        content
    };
}

function filterMetadata(metadata, metadataKeys) {
    const filtered = {};
    for (const key of metadataKeys) {
        if (metadata.hasOwnProperty(key)) {
            filtered[key] = metadata[key];
        }
    }

    return filtered;
}

(async function main() {
    await createIndex(
        "posts",
        ["id", "title", "summary", "datetime", "tags"],
        true
    );
    await createIndex("projects", ["id", "title", "short_description", "tags"]);
    await createFeed();
})().catch(console.error);
