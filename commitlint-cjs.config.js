export default {
	extends: '@commitlint/config-conventional',
	rules: {
		"header-max-length": [2, "always", 75],
		"body-max-line-length": [1, "always", 120],
		"type-min-length": [2, "always", 3],
		"type-enum": [
			2,
			"always",
			["feat", "fix", "docs", "refactor", "test", "revert", "chore", "build"],
			],
	},
};