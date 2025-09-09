// @ts-check
import withNuxt from './.nuxt/eslint.config.mjs';

import oxlint from 'eslint-plugin-oxlint';

export default withNuxt(
  // Your custom configs here
  oxlint.configs['flat/recommended'],
  {
    rules: {
      'vue/no-multiple-template-root': 'off',
      'vue/html-self-closing': [
        'error',
        {
          html: {
            void: 'always',
            normal: 'always',
            component: 'always',
          },
          svg: 'always',
          math: 'always',
        },
      ],
    },
  },
);
