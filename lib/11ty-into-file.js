const { extname } = require("path");
const CleanCSS = require("clean-css");
const { isProduction } = require("./utils/config");
const { writeStaticFile } = require("./utils/files");

module.exports = {
	configFunction(eleventyConfig) {
		const cleanCSS = new CleanCSS();

		eleventyConfig.addNunjucksAsyncFilter(
			"intoFile",
			async (data, fileName, callback) => {
				try {
					if (isProduction && extname(fileName) === ".css") {
						data = cleanCSS.minify(data).styles;
					}

					const httpPath = await writeStaticFile(fileName, data);

					callback(null, httpPath);
				} catch (err) {
					callback(err);
				}
			}
		);
	}
};
