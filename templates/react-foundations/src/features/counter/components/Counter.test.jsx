import { expect, it } from 'vitest';
import { render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Counter } from './Counter.jsx';

it('supports increment and reset through accessible controls', async () => {
  const user = userEvent.setup();
  render(<Counter />);
  await user.click(screen.getByRole('button', { name: 'Increase count' }));
  await user.click(screen.getByRole('button', { name: 'Increase count' }));
  expect(screen.getByRole('status', { name: 'Count' })).toHaveTextContent('2');
  await user.click(screen.getByRole('button', { name: 'Reset' }));
  expect(screen.getByRole('status', { name: 'Count' })).toHaveTextContent('0');
});

it('keeps separate feature instances independent', async () => {
  render(
    <>
      <Counter />
      <Counter />
    </>,
  );
  const [first, second] = screen.getAllByRole('region', {
    name: 'Example counter',
  });
  await userEvent
    .setup()
    .click(within(first).getByRole('button', { name: 'Increase count' }));
  expect(
    within(first).getByRole('status', { name: 'Count' }),
  ).toHaveTextContent('1');
  expect(
    within(second).getByRole('status', { name: 'Count' }),
  ).toHaveTextContent('0');
});
