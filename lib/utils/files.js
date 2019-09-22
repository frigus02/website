const { createHash } = require("crypto");
const path = require("path");
const { mkdir } = require("fs").promises;

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
  const ext = path.extname(name);
  const nameWithoutExt = name.substr(0, name.length - ext.length);
  return `${nameWithoutExt}-${rev}${ext}`;
};

module.exports = {
  ensureFolder,
  getNameWithHash
};
