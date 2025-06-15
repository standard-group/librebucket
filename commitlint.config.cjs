module.exports = {
    extends: [
      '@commitlint/config-conventional',
    ],
    rules: {
      'type-enum': [
        2, // rules lvls: 0-off, 1-warn, 2-error
        'always', // something
        [ // value
          'build',
          'chore',
          'ci',
          'docs',
          'feat',
          'fix',
          'perf',
          'refactor',
          'revert',
          'style',
          'test',
          'merge',
        ],
      ],
    },
  };
  