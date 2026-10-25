import { createHomePage } from '../pages/home/home-page.js';
import { createIdentityShell } from '../shared/identity/dom.js';

export function mountApp(root) {
  const page = createHomePage();
  const shell = createIdentityShell(page.element);
  root.replaceChildren(shell.element);
  return () => {
    page.destroy();
    shell.destroy();
  };
}
