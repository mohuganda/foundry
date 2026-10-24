import { createButton } from '../../shared/ui/button.js';

// Each instance owns its state and listeners; destroy it when leaving the page.
export function createCounter() {
  let count = 0;
  const controller = new AbortController();
  const element = document.createElement('section');
  element.className = 'card';
  element.setAttribute('aria-label', 'Example counter');
  const heading = document.createElement('h2');
  heading.textContent = 'Your first feature';
  const description = document.createElement('p');
  description.textContent =
    'An isolated component with accessible controls and its own tests.';
  const output = document.createElement('output');
  output.className = 'counter-value';
  output.setAttribute('aria-label', 'Count');
  output.setAttribute('aria-live', 'polite');
  output.textContent = '0';
  const actions = document.createElement('div');
  actions.className = 'actions';
  const increment = createButton('Increase count');
  const reset = createButton('Reset');
  increment.addEventListener(
    'click',
    () => {
      output.textContent = String(++count);
    },
    { signal: controller.signal },
  );
  reset.addEventListener(
    'click',
    () => {
      count = 0;
      output.textContent = '0';
    },
    { signal: controller.signal },
  );
  actions.append(increment, reset);
  element.append(heading, description, output, actions);
  return {
    element,
    destroy() {
      controller.abort();
      element.remove();
    },
  };
}
