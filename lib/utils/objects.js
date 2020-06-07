const transformKeys = (obj, fn) => {
	const newObj = {};
	for (const key of Object.keys(obj)) {
		newObj[fn(key)] = obj[key];
	}

	return newObj;
};

module.exports = {
	transformKeys,
};
