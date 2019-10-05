const { createHash } = require("crypto");
const { writeFile } = require("fs").promises;
const { promisify } = require("util");
const mkdirp = promisify(require("mkdirp"));
const { extname, dirname, basename, join: joinPath } = require("path").posix;
const { outputDir, isProduction } = require("./config");

const getNameWithHash = (name, data) => {
	const hash = createHash("md5").update(data);
	const rev = hash.digest("hex").substr(0, 8);
	const ext = extname(name);
	const nameWithoutExt = name.substr(0, name.length - ext.length);
	return `${nameWithoutExt}-${rev}${ext}`;
};

const writeStaticFile = async (filePath, data) => {
	const baseDir = joinPath(outputDir, "static");
	const dir = dirname(filePath);
	const name = basename(filePath);
	const newName = isProduction ? getNameWithHash(name, data) : name;
	await mkdirp(joinPath(baseDir, dir));
	await writeFile(joinPath(baseDir, dir, newName), data);
	return joinPath("/static", dir, newName);
};

module.exports = {
	writeStaticFile
};
