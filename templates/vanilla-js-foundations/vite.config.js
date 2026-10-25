import { defineConfig } from 'vitest/config';
import { readFileSync } from 'node:fs';
import { env as processEnv } from 'node:process';

// Set both paths to opt into local HTTPS; these values are not exposed to clients.
const https =
  processEnv.DEV_TLS_CERT && processEnv.DEV_TLS_KEY
    ? {
        cert: readFileSync(processEnv.DEV_TLS_CERT),
        key: readFileSync(processEnv.DEV_TLS_KEY),
      }
    : undefined;

export default defineConfig({
  server: { https, host: 'localhost', port: 5173, strictPort: true },
  test: {
    environment: 'happy-dom',
    setupFiles: ['./src/test/setup.js'],
    clearMocks: true,
  },
});
