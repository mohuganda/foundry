import { beforeEach, expect, it, vi } from 'vitest';
import { fireEvent, render, screen } from '@testing-library/react';
vi.mock('../shared/auth/client.js', () => ({
  auth: {
    subscribe: () => () => {},
    getSession: vi.fn(),
    initialize: vi.fn(),
    signIn: vi.fn(),
    signOut: vi.fn(),
  },
}));
import { auth } from '../shared/auth/client.js';
import { App } from './App.jsx';
beforeEach(() => {
  vi.clearAllMocks();
});
it('requires sign-in before showing authenticated content', () => {
  auth.getSession.mockReturnValue({
    status: 'anonymous',
    user: null,
    error: '',
  });
  render(<App />);
  fireEvent.click(screen.getByRole('button', { name: 'Continue' }));
  expect(auth.signIn).toHaveBeenCalledOnce();
  expect(screen.queryByRole('button', { name: 'Sign out' })).toBeNull();
});
it('exposes logout for authenticated users', () => {
  auth.getSession.mockReturnValue({
    status: 'authenticated',
    user: { subject: 'person' },
    error: '',
  });
  render(<App />);
  fireEvent.click(screen.getByRole('button', { name: 'Sign out' }));
  expect(auth.signOut).toHaveBeenCalledOnce();
  expect(screen.queryByRole('button', { name: 'Continue' })).toBeNull();
});
it('keeps the default foundation accessible when authentication is disabled', () => {
  auth.getSession.mockReturnValue({
    status: 'disabled',
    user: null,
    error: '',
  });
  render(<App />);
  expect(screen.getByRole('heading', { level: 1 })).toBeTruthy();
  expect(screen.queryByRole('button', { name: 'Continue' })).toBeNull();
});
