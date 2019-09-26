const padLeft = (padding, length, str) =>
	`${padding.repeat(length)}${str}`.slice(-length);

module.exports = {
	configFunction(eleventyConfig) {
		eleventyConfig.addFilter("formatDate", value => {
			const date = new Date(value);
			const year = date.getFullYear();
			const month = padLeft("0", 2, date.getMonth());
			const day = padLeft("0", 2, date.getDate());
			const hour = padLeft("0", 2, date.getHours());
			const minute = padLeft("0", 2, date.getMinutes());
			return `${year}-${month}-${day} ${hour}:${minute}`;
		});
	}
};
