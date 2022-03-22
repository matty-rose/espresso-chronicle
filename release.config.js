let scope = process.env.SCOPE;
let tagFormat = scope + "-v${version}"
let headerPatternRegex = new RegExp(
  "^(\\w*)(?:\\((" + scope + ")\\))?\\: (.*)$"
);

let plugins = [
  [
    "@semantic-release/commit-analyzer",
    {
      preset: "conventionalcommits",
      releaseRules: [{ scope: `!${scope}`, release: false }],
    },
  ],
];

plugins.push([
  "@semantic-release/exec",
  {
    verifyReleaseCmd: "echo ${nextRelease.version} > .version"
  }
])

plugins.push([
  "@semantic-release/release-notes-generator",
  {
    preset: "conventionalcommits",
    parserOpts: {
      headerPattern: headerPatternRegex,
    },
  },
]);

plugins.push([
  "@semantic-release/exec",
  {
    prepareCmd: 'git fetch origin "refs/tags/*:refs/tags/*"'
  }
])

plugins.push("@semantic-release/github");

module.exports = {
  branches: ["main"],
  tagFormat: tagFormat,
  plugins: plugins,
};
