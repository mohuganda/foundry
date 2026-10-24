import { expect, it } from 'vitest';
import { getByRole } from '@testing-library/dom';
import userEvent from '@testing-library/user-event';
import { createCounter } from './counter.js';

it('supports increment and reset through accessible controls', async () => {
  const counter = createCounter();
  document.body.append(counter.element);
  const user = userEvent.setup();
  await user.click(
    getByRole(counter.element, 'button', { name: 'Increase count' }),
  );
  await user.click(
    getByRole(counter.element, 'button', { name: 'Increase count' }),
  );
  expect(
    getByRole(counter.element, 'status', { name: 'Count' }),
  ).toHaveTextContent('2');
  await user.click(getByRole(counter.element, 'button', { name: 'Reset' }));
  expect(
    getByRole(counter.element, 'status', { name: 'Count' }),
  ).toHaveTextContent('0');
  counter.destroy();
});

it('keeps instances independent and removes listeners on teardown', async () => {
  const first = createCounter();
  const second = createCounter();
  document.body.append(first.element, second.element);
  const button = getByRole(first.element, 'button', { name: 'Increase count' });
  const output = getByRole(first.element, 'status', { name: 'Count' });
  await userEvent.setup().click(button);
  expect(
    getByRole(second.element, 'status', { name: 'Count' }),
  ).toHaveTextContent('0');
  first.destroy();
  button.click();
  expect(output).toHaveTextContent('1');
  expect(first.element).not.toBeInTheDocument();
  second.destroy();
});
