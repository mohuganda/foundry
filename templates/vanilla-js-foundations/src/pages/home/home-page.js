import { createCounter } from '../../features/counter/counter.js';
import { identity } from '../../shared/identity/index.js';

export function createHomePage() {
  const element = document.createElement('div');
  const hero = document.createElement('section');
  hero.className = 'hero';
  const heading = document.createElement('h1');
  heading.textContent = identity.application.name;
  const description = document.createElement('p');
  description.className = 'muted';
  description.textContent = identity.application.description;
  hero.append(heading, description);
  element.append(hero);
  const counter = createCounter();
  element.append(counter.element);
  return {
    element,
    destroy() {
      counter.destroy();
      element.remove();
    },
  };
}
