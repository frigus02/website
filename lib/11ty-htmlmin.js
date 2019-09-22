const htmlmin = require("html-minifier");

module.exports = {
  configFunction(eleventyConfig) {
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
