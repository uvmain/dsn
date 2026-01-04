import antfu from '@antfu/eslint-config'

export default antfu({
  vue: true,
  unocss: true,
  typescript: {
    tsconfigPath: 'tsconfig.json',
  },
  ignores: ['dist', '**/dist/**', 'auto-imports.d.ts', 'components.d.ts'],
}, {
  files: ['**/*.{vue,ts,js,json}'],
  rules: {
    'no-console': 'off',
    'brace-style': ['error', 'stroustrup'],
    'curly': ['off'],
    'vue/html-self-closing': 'off',
    '@typescript-eslint/no-explicit-any': 'off',
    '@typescript-eslint/no-unsafe-argument': 'off',
    '@typescript-eslint/no-unsafe-assignment': 'off',
    '@typescript-eslint/no-unsafe-call': 'off',
    '@typescript-eslint/no-unsafe-return': 'off',
    'import/order': 'off',
  },
})
