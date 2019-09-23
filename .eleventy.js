const pluginSyntaxHighlight = require("@11ty/eleventy-plugin-syntaxhighlight");
const pluginRss = require("@11ty/eleventy-plugin-rss");
const pluginFormatDate = require("./lib/11ty-format-date");
const pluginHtmlMin = require("./lib/11ty-htmlmin");
const pluginImage = require("./lib/11ty-image");
const pluginIntoFile = require("./lib/11ty-into-file");
const pluginSubHeading = require("./lib/11ty-sub-heading");

module.exports = eleventyConfig => {
	eleventyConfig.addPlugin(pluginFormatDate);
	eleventyConfig.addPlugin(pluginHtmlMin);
	eleventyConfig.addPlugin(pluginImage);
	eleventyConfig.addPlugin(pluginIntoFile);
	eleventyConfig.addPlugin(pluginSubHeading);
	eleventyConfig.addPlugin(pluginSyntaxHighlight);
	eleventyConfig.addPlugin(pluginRss);

	eleventyConfig.addPassthroughCopy("**/.htaccess");
	eleventyConfig.addPassthroughCopy("src/favicon.ico");

	return {
		dir: {
			input: "src",
			output: "build"
		},
		htmlTemplateEngine: "njk"
	};
};
