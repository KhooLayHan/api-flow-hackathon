// @ts-check
import withNuxt from './.nuxt/eslint.config.mjs';

import oxlint from 'eslint-plugin-oxlint';

export default withNuxt(
  // Your custom configs here
  oxlint.configs['flat/recommended'],
);
