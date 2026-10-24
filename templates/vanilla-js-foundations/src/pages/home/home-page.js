import { createCounter } from '../../features/counter/counter.js';
import { env } from '../../shared/config/env.js';

export function createHomePage() {
  const element = document.createElement('div');
  element.className = 'shell';
  element.innerHTML = `
    <a class="skip-link" href="#main">Skip to content</a>
    <header class="brand">Foundry / JavaScript</header>
    <main id="main" tabindex="-1">
      <section class="hero" aria-labelledby="page-title">
        <p class="eyebrow">A foundation for your next idea</p>
        <h1 id="page-title"></h1>
        <p class="muted">Native JavaScript, a clear feature structure, and room to grow.</p>
      </section>
    </main>
    <footer class="muted">Start with one feature. Build from there.</footer>
  `;
  element.querySelector('h1').textContent = env.appName;
  const counter = createCounter();
  element.querySelector('main').append(counter.element);
  return {
    element,
    destroy() {
      counter.destroy();
      element.remove();
    },
  };
}
