import { rmSync } from 'node:fs';

// Resolve from this file so only this project's generated artifacts are removed.
const projectRoot = new URL('../', import.meta.url);
for (const directory of ['dist/', 'coverage/', '.vitest/']) {
  rmSync(new URL(directory, projectRoot), { recursive: true, force: true });
}
