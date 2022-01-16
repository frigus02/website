const { join: joinPath } = require("path").posix;
const pluginSyntaxHighlight = require("@11ty/eleventy-plugin-syntaxhighlight");
const pluginRss = require("@11ty/eleventy-plugin-rss");
const markdownIt = require("markdown-it");
const markdownItAnchor = require("markdown-it-anchor");
const pluginFormatDate = require("./lib/11ty-format-date");
const pluginHtmlImages = require("./lib/11ty-html-images");
const pluginHtmlMin = require("./lib/11ty-htmlmin");
const pluginIntoFile = require("./lib/11ty-into-file");
const pluginSubHeading = require("./lib/11ty-sub-heading");
const { inputDir, outputDir } = require("./lib/utils/config");

module.exports = (eleventyConfig) => {
	eleventyConfig.addPlugin(pluginFormatDate);
	eleventyConfig.addPlugin(pluginHtmlImages);
	eleventyConfig.addPlugin(pluginHtmlMin);
	eleventyConfig.addPlugin(pluginIntoFile);
	eleventyConfig.addPlugin(pluginSubHeading);
	eleventyConfig.addPlugin(pluginSyntaxHighlight);
	eleventyConfig.addPlugin(pluginRss);

	eleventyConfig.addPassthroughCopy(joinPath(inputDir, "favicon.ico"));

	const markdownLib = markdownIt({ html: true }).use(markdownItAnchor);
	eleventyConfig.setLibrary("md", markdownLib);

	return {
		dir: {
			input: inputDir,
			output: outputDir,
		},
		htmlTemplateEngine: "njk",
		markdownTemplateEngine: false,
	};
};
