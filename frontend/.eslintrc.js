module.exports = {
  extends: require.resolve('@umijs/max/eslint'),
  rules: {
    'comma-dangle': [
      'error',
      {
        arrays: 'only-multiline',
        objects: 'only-multiline',
        imports: 'only-multiline',
        exports: 'only-multiline',
        functions: 'always-multiline',
      },
    ],
    eqeqeq: 'off',
    'no-unused-expressions': 'off',
    '@typescript-eslint/no-unused-expressions': ['off'],
  },
};
