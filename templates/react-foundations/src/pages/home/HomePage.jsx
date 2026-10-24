import { Counter } from '../../features/counter/components/Counter.jsx';
import { env } from '../../shared/config/env.js';

export function HomePage() {
  return (
    <div className="shell">
      <a className="skip-link" href="#main">
        Skip to content
      </a>
      <header className="brand">Foundry / React</header>
      <main id="main" tabIndex={-1}>
        <section className="hero" aria-labelledby="page-title">
          <p className="eyebrow">A foundation for your next idea</p>
          <h1 id="page-title">{env.appName}</h1>
          <p className="muted">
            React, a clear feature structure, and room to grow.
          </p>
        </section>
        <Counter />
      </main>
      <footer className="muted">
        Start with one feature. Build from there.
      </footer>
    </div>
  );
}
