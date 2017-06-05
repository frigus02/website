/* eslint-env node */

'use strict';

const fs = require('fs-extra');
const readline = require('readline');

async function createIndex(folder, metadataKeys, reverse) {
    await fs.remove(`data-api/${folder}`);
    await fs.mkdirs(`data-api/${folder}`);

    const files = await fs.readdir(`data/${folder}`);
    let index = [];

    for (const fileName of files) {
        const {metadata} = await getMetadata(`data/${folder}/${fileName}`);
        metadata.id = fileName.split('-')[1].split('.')[0];
        index.push(metadata);

        await fs.copy(`data/${folder}/${fileName}`, `data-api/${folder}/${metadata.id}.md`);
    }

    index = index.map(metadata => filterMetadata(metadata, metadataKeys));
    if (reverse) {
        index.reverse();
    }

    const indexJson = JSON.stringify(index);
    await fs.writeFile(`data-api/${folder}/index.json`, indexJson, 'utf-8');
}

async function getMetadata(file) {
    const text = await fs.readFile(file, 'utf-8');

    const metadataStart = text.indexOf('```json') + 7;
    const metadataEnd = text.indexOf('```', metadataStart);
    const metadata = text.substring(metadataStart, metadataEnd);

    const content = text.substring(metadataEnd + 3);

    return {
        metadata: JSON.parse(metadata),
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

    await createIndex('posts', ['id', 'title', 'summary', 'datetime', 'tags'], true);
    await createIndex('projects', ['id', 'title', 'short_description', 'tags']);

})().catch(console.error);
