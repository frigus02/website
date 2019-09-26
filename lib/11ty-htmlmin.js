const htmlmin = require("html-minifier");
const { isProduction } = require("./utils/config");

module.exports = {
	configFunction(eleventyConfig) {
		if (!isProduction) return;

		eleventyConfig.addTransform("htmlmin", (content, outputPath) => {
			if (outputPath.endsWith(".html")) {
				return htmlmin.minify(content, {
					collapseWhitespace: true,
					removeComments: true
				});
			}

			return content;
		});
	}
};
