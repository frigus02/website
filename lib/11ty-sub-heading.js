const { resolve: resolvePath } = require("path");
const nunjucks = require("nunjucks");
const { inputDir } = require("./utils/config");

module.exports = {
	configFunction(eleventyConfig) {
		eleventyConfig.addShortcode("subHeading", (page, title) => {
			if (page.url === "/") return "";

			const url = page.url.substr(0, page.url.indexOf("/", 1) + 1);
			return new nunjucks.runtime.SafeString(
				resolvePath(page.inputPath) === resolvePath(inputDir, "404.html")
					? `<h2>${title}</h2>`
					: `<h2><a href="${url}">${title}</a></h2>`
			);
		});
	}
};
