const path = require("path");
const { writeFile } = require("fs").promises;
const CleanCSS = require("clean-css");
const { getNameWithHash, ensureFolder } = require("./utils/files");

module.exports = {
  configFunction(eleventyConfig) {
    const cleanCSS = new CleanCSS();

    eleventyConfig.addNunjucksAsyncFilter(
      "intoFile",
      async (data, fileName, callback) => {
        try {
          if (path.extname(fileName) === ".css") {
            data = cleanCSS.minify(data).styles;
          }

          const name = getNameWithHash(fileName, data);
          const outputPath = path.resolve("./build/static/", name);
          await ensureFolder(path.dirname(outputPath));
          await writeFile(outputPath, data);

          callback(null, `/static/${name}`);
        } catch (err) {
          callback(err);
        }
      }
    );
  }
};
