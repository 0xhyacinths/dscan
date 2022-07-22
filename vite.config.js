import { sveltekit } from '@sveltejs/kit/vite';
import * as child from 'child_process';

const revision = child.execSync('git rev-parse HEAD').toString().trim()

/** @type {import('vite').UserConfig} */
const config = {
	plugins: [sveltekit()],
  define: {
    '__APP_VERSION__': JSON.stringify(revision.slice(0,8)),
  }
};

export default config;
