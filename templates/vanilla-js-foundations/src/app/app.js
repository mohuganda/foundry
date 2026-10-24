import { createHomePage } from '../pages/home/home-page.js';

export function mountApp(root) {
  const page = createHomePage();
  root.replaceChildren(page.element);
  return () => page.destroy();
}
