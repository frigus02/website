class Cache {
	constructor() {
		this.cache = new Map();
	}

	get(key, action) {
		let value = this.cache.get(key);
		if (!value) {
			value = action();
			this.cache.set(key, value);
		}

		return value;
	}
}

module.exports = Cache;
