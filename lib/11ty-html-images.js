const { readFile } = require("fs").promises;
const {
	join: joinPath,
	dirname,
	relative: relativePath,
	resolve: resolvePath
} = require("path");
const posthtml = require("posthtml");
const Cache = require("./utils/cache");
const { inputDir, outputDir } = require("./utils/config");
const { writeStaticFile } = require("./utils/files");

const plugin = opts => tree => {
	const promises = [];

	tree.match({ tag: "img" }, node => {
		const src = node.attrs && node.attrs.src;
		if (!src) return node;

		// Skip external images
		if (/^https?:\/\//.test(src)) return node;

		const inputPath = src.startsWith("/")
			? joinPath(inputDir, src)
			: joinPath(opts.cwd, src);

		// Root path relative to inputDir
		//   src/static/images/me.png
		//     outputRelativePath me.png
		//     outputRootedPath   /me.png
		//   src/projects/rester/preview.png
		//     outputRelativePath ../../projects/rester/preview.png
		//     outputRootedPath   /projects/rester/preview.png
		const outputRelativePath = relativePath(
			joinPath(inputDir, "static/images"),
			inputPath
		);
		const outputRootedPath = resolvePath("/", outputRelativePath);
		const outputPath = joinPath("images", outputRootedPath);

		promises.push(
			(async () => {
				const newSrc = await opts.cache.get(inputPath, async () => {
					try {
						const data = await readFile(inputPath);
						return await writeStaticFile(outputPath, data);
					} catch (err) {
						if (err.code === "ENOENT") {
							console.warn(`Image ${inputPath} does not exist`);
							return "";
						}
					}
				});
				node.attrs.src = newSrc;
			})()
		);

		return node;
	});

	return Promise.all(promises)
		.catch(err => {
			console.error(err);
		})
		.then(() => tree);
};

module.exports = {
	configFunction(eleventyConfig) {
		const cache = new Cache();
		eleventyConfig.addTransform("images", async (content, outputPath) => {
			if (outputPath.endsWith(".html")) {
				const cwd = joinPath(
					inputDir,
					relativePath(outputDir, dirname(outputPath))
				);
				const result = await posthtml()
					.use(plugin({ cache, cwd }))
					.process(content);
				return result.html;
			}

			return content;
		});
	}
};
