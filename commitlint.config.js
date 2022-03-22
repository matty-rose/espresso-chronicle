module.exports = {
  extends: ["@commitlint/config-conventional"],
  rules: {
    "header-max-length": [2, "always", 100],
    "references-empty": [0, "never"],
    "scope-case": [2, "always", "lower-case"],
    "scope-empty": [2, "never"],
    "subject-case": [0],
    "type-case": [2, "always", "lower-case"],
    "scope-enum": [2, "always", ["repo", "api", "infra"]]
  },
};
