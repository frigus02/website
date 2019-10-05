const { readFile, unlink, writeFile } = require("fs").promises;
const { join: joinPath, relative: relativePath } = require("path");
const chokidar = require("chokidar");
const CleanCSS = require("clean-css");
const pLimit = require("p-limit");
const { inputDir, isProduction, outputDir } = require("./utils/config");
const { writeStaticFile } = require("./utils/files");
const { transformKeys } = require("./utils/objects");

class Input {
	constructor() {
		this.files = {};
	}

	async update(event, path) {
		switch (event) {
			case "add":
			case "change":
				this.files[path] = await readFile(path, "utf8");
				break;
			case "unlink":
				delete this.files[path];
				break;
		}

		return this.files;
	}
}

class Transform {
	constructor() {
		this.cleanCSS = new CleanCSS();
		this.concatOutputPath = "styles.css";
	}

	async update(files) {
		return isProduction
			? { [this.concatOutputPath]: this.concatAndMinify(files) }
			: transformKeys(files, path => relativePath(`${inputDir}/static`, path));
	}

	concatAndMinify(inputFiles) {
		const concat = Object.keys(inputFiles)
			.sort()
			.map(path => inputFiles[path])
			.join("\n\n");
		return this.cleanCSS.minify(concat).styles;
	}
}

class Output {
	constructor() {
		this.dataOutputPath = "_data/styles.json";
		this.writtenFiles = [];
	}

	async update(files) {
		const newFiles = await Promise.all(
			Object.keys(files).map(path => writeStaticFile(path, files[path]))
		);
		const removedFiles = this.writtenFiles.filter(
			file => !newFiles.includes(file)
		);
		for (const file of removedFiles) {
			await unlink(joinPath(outputDir, file));
		}

		this.writtenFiles = newFiles.sort();
		await this.writeDataFile();
	}

	async writeDataFile() {
		await writeFile(
			joinPath(inputDir, this.dataOutputPath),
			JSON.stringify(this.writtenFiles),
			"utf8"
		);
	}
}

const main = async () => {
	const args = process.argv.slice(2);
	const watch = args[0] === "--watch";

	const input = new Input();
	const transform = new Transform();
	const output = new Output();
	const update = async (event, path) =>
		await output.update(
			await transform.update(await input.update(event, path))
		);

	const fsLimit = pLimit(1);
	const watcher = chokidar.watch(`${inputDir}/**/*.css`, {
		persistent: watch
	});
	watcher.on("all", (event, path) => {
		console.log(event, path);
		fsLimit(update, event, path);
	});
};

main().catch(err => {
	console.error(err);
	process.exitCode = 1;
});
