const { readFile } = require("fs").promises;
const { join: joinPath, dirname } = require("path");
const { inputDir } = require("./utils/config");
const { writeStaticFile } = require("./utils/files");

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

const isPostOrProject = page =>
	page && /\/(projects|posts)\/[a-z0-9_]+\//.test(page.url);

const getInputOutputPaths = (name, page) =>
	isPostOrProject(page)
		? {
				in: joinPath(dirname(page.inputPath), name),
				out: joinPath(page.url, name)
		  }
		: {
				in: joinPath(inputDir, "static/images", name),
				out: joinPath("images", name)
		  };

module.exports = {
	configFunction(eleventyConfig) {
		const cache = new ImageCache();

		eleventyConfig.addNunjucksAsyncFilter(
			"image",
			async (name, page, callback) => {
				try {
					if (typeof callback === "undefined") {
						callback = page;
						page = undefined;
					}

					const paths = getInputOutputPaths(name, page);
					const httpPath = await cache.get(paths.in, async () => {
						try {
							const data = await readFile(paths.in);
							return await writeStaticFile(paths.out, data);
						} catch (err) {
							if (err.code === "ENOENT") {
								console.warn(`Image ${paths.in} does not exist`);
								return "";
							}
						}
					});

					callback(null, httpPath);
				} catch (err) {
					callback(err);
				}
			}
		);
	}
};
