const { createHash } = require("crypto");
const { mkdir, writeFile } = require("fs").promises;
const { extname, dirname, basename, join: joinPath } = require("path");
const { outputDir, isProduction } = require("./config");

const ensureFolder = async path => {
	try {
		await mkdir(path, { recursive: true });
	} catch (err) {
		if (err.code !== "EEXIST") {
			throw err;
		}
	}
};

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
	await ensureFolder(joinPath(baseDir, dir));
	await writeFile(joinPath(baseDir, dir, newName), data);
	return joinPath("/static", dir, newName);
};

module.exports = {
	writeStaticFile
};
