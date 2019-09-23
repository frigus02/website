const path = require("path");
const { readFile, writeFile } = require("fs").promises;
const { getNameWithHash, ensureFolder } = require("./utils/files");

class ImageCache {
	constructor() {
		this.cache = new Map();
	}

	get(key, action) {
		let value = this.cache.get(key);
		if (!value) {
			value = action();
			this.cache.set(key, value);
		}

		return value;
	}
}

module.exports = {
	configFunction(eleventyConfig) {
		const cache = new ImageCache();
		const specialUrl = /\/(projects|posts)\/[a-z0-9_]+\//;

		eleventyConfig.addNunjucksAsyncFilter(
			"image",
			async (name, page, callback) => {
				try {
					if (typeof callback === "undefined") {
						callback = page;
						page = undefined;
					}

					let inputDirectory = "./src/static/images/";
					let outputDirectory = "./build/static/images/";
					let httpDirectory = "/static/images/";
					if (page && specialUrl.test(page.url)) {
						inputDirectory = path.dirname(page.inputPath);
						outputDirectory = `./build/static${page.url}`;
						httpDirectory = `/static${page.url}`;
					}

					const inputPath = path.resolve(inputDirectory, name);
					const httpPath = await cache.get(inputPath, async () => {
						const data = await readFile(inputPath);
						const newName = getNameWithHash(name, data);

						const outputPath = path.resolve(outputDirectory, newName);
						await ensureFolder(path.dirname(outputPath));
						await writeFile(outputPath, data);

						return `${httpDirectory}${newName}`;
					});

					callback(null, httpPath);
				} catch (err) {
					callback(err);
				}
			}
		);
	}
};
