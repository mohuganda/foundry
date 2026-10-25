// Visual integration check against the disposable local Keycloak fixture.
import assert from 'node:assert/strict';
import { createRequire } from 'node:module';
import { mkdir, writeFile } from 'node:fs/promises';
import process from 'node:process';

if (process.argv.includes('--help') || process.argv.includes('-h')) {
  console.log('Usage: node scripts/check-login-theme.mjs');
  console.log(
    'Requires the local Keycloak fixture, FOUNDRY_PLAYWRIGHT_PACKAGE (absolute package.json), FOUNDRY_CHROMIUM_EXECUTABLE, and FOUNDRY_LOGIN_SCREENSHOTS. Captures desktop/light and mobile/dark previews without entering credentials.',
  );
  process.exit(0);
}
assert.equal(process.argv.length, 2, 'Unexpected arguments; use --help.');
for (const key of [
  'FOUNDRY_PLAYWRIGHT_PACKAGE',
  'FOUNDRY_CHROMIUM_EXECUTABLE',
  'FOUNDRY_LOGIN_SCREENSHOTS',
])
  assert(process.env[key], `${key} is required`);

const require = createRequire(process.env.FOUNDRY_PLAYWRIGHT_PACKAGE);
const { chromium } = require('playwright');
const output = process.env.FOUNDRY_LOGIN_SCREENSHOTS;
await mkdir(output, { recursive: true });
const browser = await chromium.launch({
  executablePath: process.env.FOUNDRY_CHROMIUM_EXECUTABLE,
  headless: true,
});
try {
  for (const [name, viewport, colorScheme] of [
    ['desktop-light', { width: 1280, height: 900 }, 'light'],
    ['mobile-dark', { width: 375, height: 812 }, 'dark'],
  ]) {
    // Only this disposable loopback provider uses a self-signed certificate.
    const context = await browser.newContext({
      viewport,
      colorScheme,
      ignoreHTTPSErrors: true,
    });
    const page = await context.newPage();
    const params = new URLSearchParams({
      client_id: 'programme-web',
      redirect_uri: 'https://localhost:5173/auth/callback',
      response_type: 'code',
      scope: 'openid',
      state: 'theme-preview',
      nonce: 'theme-preview',
      code_challenge_method: 'S256',
      code_challenge: '012345678901234567890123456789012345678901234',
    });
    await page.goto(
      `https://localhost:8443/realms/foundry/protocol/openid-connect/auth?${params}`,
    );
    await page.locator('#username').waitFor();
    await page.evaluate(() => document.fonts.ready);
    assert.equal(await page.locator('.foundry-organisation').count(), 1);
    const metrics = await page.evaluate(() => {
      const header = document
        .querySelector('#kc-header')
        .getBoundingClientRect();
      const main = document.querySelector('main').getBoundingClientRect();
      const action = document
        .querySelector('#kc-login')
        .getBoundingClientRect();
      return {
        headerWidth: header.width,
        mainWidth: main.width,
        left: main.left,
        right: main.right,
        actionHeight: action.height,
        pageWidth: document.documentElement.scrollWidth,
        viewport: innerWidth,
      };
    });
    console.log(name, metrics);
    await page.locator('#username').focus();
    await page.screenshot({
      path: `${output}/keycloak-${name}.png`,
      fullPage: true,
    });
    const snapshot = await page.evaluate(() => {
      const copy = document.documentElement.cloneNode(true);
      for (const script of copy.querySelectorAll('script')) script.remove();
      const sheets = [...document.styleSheets].map((sheet) =>
        [...sheet.cssRules].map((rule) => rule.cssText).join('\n'),
      );
      for (const link of copy.querySelectorAll('link[rel="stylesheet"]'))
        link.remove();
      for (const css of sheets) {
        const style = document.createElement('style');
        style.textContent = css;
        copy.querySelector('head').append(style);
      }
      for (const form of copy.querySelectorAll('form'))
        form.removeAttribute('action');
      return '<!doctype html>' + copy.outerHTML;
    });
    await writeFile(`${output}/keycloak-${name}.html`, snapshot);
    assert(
      Math.abs(metrics.headerWidth - metrics.mainWidth) < 1,
      'provider card sections do not align',
    );
    assert(
      Math.abs(metrics.mainWidth - Math.min(440, viewport.width - 32)) < 1,
      'provider card width differs from shared design',
    );
    assert(
      metrics.actionHeight >= 48,
      'provider action is smaller than shared control height',
    );
    assert(
      metrics.pageWidth <= metrics.viewport &&
        metrics.left >= 0 &&
        metrics.right <= metrics.viewport,
      'provider login overflows viewport',
    );
    await context.close();
  }
} finally {
  await browser.close();
}
