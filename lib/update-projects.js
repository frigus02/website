const { readdir, readFile, stat, writeFile } = require("fs").promises;
const { extname } = require("path");
const fetch = require("node-fetch");

const forEachSource = async (path, cb) => {
	const entries = await readdir(path);
	for (const entry of entries) {
		const fullName = `${path}/${entry}`;
		const stats = await stat(fullName);
		if (stats.isFile() && entry === ".source.json") {
			const source = await readFile(fullName, "utf8");
			await cb(path, JSON.parse(source));
		} else if (stats.isDirectory()) {
			await forEachSource(fullName, cb);
		}
	}
};

const fetchBody = async (url, options) => {
	const res = await fetch(url, options);
	if (res.ok) {
		if (res.headers.get("content-type").includes("json")) {
			return res.json();
		} else {
			return res.text();
		}
	} else {
		throw new Error(`fetch ${url} -> ${res.status}: ${await res.text()}`);
	}
};

const github = async (source) => {
	const repoUrl = `https://api.github.com/repos/${source.repo}`;
	const repo = await fetchBody(repoUrl);
	const readme = await fetchBody(`${repoUrl}/readme`);
	const topics = await fetchBody(`${repoUrl}/topics`, {
		headers: {
			accept: "application/vnd.github.mercy-preview+json",
		},
	});
	const data = {
		date: repo.pushed_at,
		metadata: {
			title: extractTitleFromGitHubReadme(readme) || repo.name,
			short_description: repo.description,
			tags: topics.names,
			homepage: repo.homepage,
			language: repo.language,
			license: repo.license ? repo.license.spdx_id : null,
			source: `https://github.com/${source.repo}`,
		},
	};
	const {
		readme: preparedReadme,
		images,
	} = await extractImagesFromGitHubReadme(prepareGitHubReadme(readme));

	return {
		readme: preparedReadme,
		data,
		images,
	};
};

const extractTitleFromGitHubReadme = (readme) => {
	const title = /^# (.+)$/gm.exec(Buffer.from(readme.content, readme.encoding));
	return title && title[1];
};

const prepareGitHubReadme = (readme) =>
	Buffer.from(readme.content, readme.encoding)
		.toString()
		// Remove main heading
		.replace(/^# .+$/gm, "")
		// Make relative URLs absolute
		.replace(
			/!\[([^\[\]]*)\]\(([^\)]+)\)/g,
			(_, p1, p2) => `![${p1}](${new URL(p2, readme.download_url)})`
		)
		.replace(
			/\[([^\[\]]*)\]\((?!#)([^\)]+)\)/g,
			(_, p1, p2) => `[${p1}](${new URL(p2, readme.html_url)})`
		)
		// Replace language tags with PrismJS supported tags
		// https://prismjs.com/#supported-languages
		.replace(/```sh/g, "```shell-session")
		.replace(/```command/g, "```shell-session");

const extractImagesFromGitHubReadme = async (readme) => {
	const placeholders = {};
	readme = readme.replace(/!\[([^\[\]]*)\]\(([^\)]+)\)/g, (_, p1, p2) => {
		if (!placeholders[p2]) {
			const placeholder = `@@IMAGE-${Object.keys(placeholders).length}@@`;
			placeholders[p2] = placeholder;
		}

		return `![${p1}](${placeholders[p2]})`;
	});

	const urls = Object.keys(placeholders);
	const images = {};
	for (const [i, url] of urls.entries()) {
		const res = await fetch(url);
		if (res.ok) {
			const data = await res.buffer();
			const ext = extname(new URL(url).pathname);
			const name = `image-${i}${ext}`;
			images[name] = data;
			readme = readme.replace(placeholders[url], name);
		} else {
			console.warn(`Failed to download image ${url}`);
		}
	}

	return { readme, images };
};

const sources = {
	github,
};

const main = async () => {
	await forEachSource("src", async (path, source) => {
		if (sources[source.source]) {
			console.log(`Update project from ${source.source} in ${path}`);
			const { readme, data, images } = await sources[source.source](source);
			await writeFile(`${path}/index.md`, readme, "utf8");
			await writeFile(
				`${path}/index.11tydata.json`,
				JSON.stringify(data, null, "\t"),
				"utf8"
			);
			for (const name of Object.keys(images)) {
				await writeFile(`${path}/${name}`, images[name]);
			}
		} else {
			console.warn(`Unknown source ${source.source} in ${path}`);
		}
	});
};

main().catch((err) => {
	console.error(err);
	process.exitCode = 1;
});
