const { readdir, readFile, writeFile } = require("fs").promises;
const fetch = require("node-fetch");

const forEachSource = async (path, cb) => {
	const entries = await readdir(path, {
		withFileTypes: true
	});
	for (const entry of entries) {
		if (entry.isFile() && entry.name === ".source.json") {
			const source = await readFile(`${path}/${entry.name}`, "utf8");
			await cb(path, JSON.parse(source));
		} else if (entry.isDirectory()) {
			await forEachSource(`${path}/${entry.name}`, cb);
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

const github = async source => {
	const repoUrl = `https://api.github.com/repos/${source.repo}`;
	const repo = await fetchBody(repoUrl);
	const readme = await fetchBody(`${repoUrl}/readme`);
	const topics = await fetchBody(`${repoUrl}/topics`, {
		headers: {
			accept: "application/vnd.github.mercy-preview+json"
		}
	});
	const data = {
		date: repo.pushed_at,
		metadata: {
			title: repo.name,
			short_description: repo.description,
			tags: topics.names,
			homepage: repo.homepage,
			language: repo.language,
			license: repo.license.spdx_id
		}
	};

	return {
		readme: prepareGitHubReadme(readme),
		data
	};
};

const prepareGitHubReadme = readme =>
	Buffer.from(readme.content, readme.encoding)
		.toString()
		// Remove main heading
		.replace(/^# .+$/gm, "")
		// Make relative URLs absolute
		.replace(
			/\[([^\[\]]+)\]\(([^\)]+)\)/g,
			(_, p1, p2) => `[${p1}](${new URL(p2, readme.html_url)})`
		)
		// Replace language tags with PrismJS supported tags
		// https://prismjs.com/#supported-languages
		.replace(/```sh/g, "```shell-session")
		.replace(/```command/g, "```shell-session");

const sources = {
	github
};

const main = async () => {
	await forEachSource("src", async (path, source) => {
		if (sources[source.source]) {
			const { readme, data } = await sources[source.source](source);
			await writeFile(`${path}/index.md`, readme, "utf8");
			await writeFile(
				`${path}/index.11tydata.json`,
				JSON.stringify(data, null, "\t"),
				"utf8"
			);
		} else {
			console.warn(`Unknown source ${source.source} in ${path}`);
		}
	});
};

main().catch(err => {
	console.error(err);
	process.exitCode = 1;
});
