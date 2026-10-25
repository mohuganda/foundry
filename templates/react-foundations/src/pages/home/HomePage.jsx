import { Counter } from '../../features/counter/components/Counter.jsx';
import { identity } from '../../shared/identity/index.js';

export function HomePage() {
  return (
    <>
      <section className="hero" aria-labelledby="page-title">
        <h1 id="page-title">{identity.application.name}</h1>
        <p className="muted">{identity.application.description}</p>
      </section>
      <Counter />
    </>
  );
}
