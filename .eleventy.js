const { join: joinPath } = require("path");
const pluginSyntaxHighlight = require("@11ty/eleventy-plugin-syntaxhighlight");
const pluginRss = require("@11ty/eleventy-plugin-rss");
const pluginFormatDate = require("./lib/11ty-format-date");
const pluginHtmlMin = require("./lib/11ty-htmlmin");
const pluginImage = require("./lib/11ty-image");
const pluginIntoFile = require("./lib/11ty-into-file");
const pluginSubHeading = require("./lib/11ty-sub-heading");
const { inputDir, outputDir } = require("./lib/utils/config");

module.exports = eleventyConfig => {
	eleventyConfig.addPlugin(pluginFormatDate);
	eleventyConfig.addPlugin(pluginHtmlMin);
	eleventyConfig.addPlugin(pluginImage);
	eleventyConfig.addPlugin(pluginIntoFile);
	eleventyConfig.addPlugin(pluginSubHeading);
	eleventyConfig.addPlugin(pluginSyntaxHighlight);
	eleventyConfig.addPlugin(pluginRss);

	eleventyConfig.addPassthroughCopy(joinPath(inputDir, "**/.htaccess"));
	eleventyConfig.addPassthroughCopy(joinPath(inputDir, "favicon.ico"));

	return {
		dir: {
			input: inputDir,
			output: outputDir
		},
		htmlTemplateEngine: "njk"
	};
};
