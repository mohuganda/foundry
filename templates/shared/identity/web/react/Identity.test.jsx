import { expect, it, vi } from 'vitest';
import { fireEvent, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { identity } from './index.js';
import {
  IdentityShell,
  IdentitySignIn,
  IdentityLoading,
  IdentityError,
  IdentityReportHeader,
} from './Identity.jsx';

it('supports missing or broken logos while retaining organisation identity', () => {
  const profile = structuredClone(identity);
  profile.organisation.name = 'Ministry of Health';
  profile.branding.assets.logo_light = '/branding/assets/logo_light.svg';
  profile.branding.assets.coat_of_arms = '';
  profile.application.environment = 'production';
  const { container } = render(
    <IdentityShell profile={profile}>
      <h1>Dashboard</h1>
    </IdentityShell>,
  );
  fireEvent.error(screen.getByRole('img', { name: 'Ministry of Health logo' }));
  expect(screen.queryByRole('img')).toBeNull();
  expect(screen.getByRole('banner')).toHaveTextContent('Ministry of Health');
  expect(screen.getByRole('main')).toContainElement(
    screen.getByRole('heading', { name: 'Dashboard' }),
  );
  expect(container.querySelector('.identity-environment')).toBeNull();
});

it('connects sign-in and retry to application callbacks', async () => {
  const signIn = vi.fn();
  const retry = vi.fn();
  render(
    <>
      <IdentitySignIn onSignIn={signIn} />
      <IdentityError message="Offline" onRetry={retry} />
    </>,
  );
  await userEvent
    .setup()
    .click(screen.getByRole('button', { name: 'Continue' }));
  await userEvent
    .setup()
    .click(screen.getByRole('button', { name: 'Try again' }));
  expect(signIn).toHaveBeenCalledTimes(1);
  expect(retry).toHaveBeenCalledTimes(1);
  expect(screen.getByRole('alert')).toHaveTextContent('Offline');
});

it('provides loading status and explicit report metadata', () => {
  render(
    <>
      <IdentityLoading message="Loading records" />
      <IdentityReportHeader title="Summary" version="3" issuedOn="2026-10-25" />
    </>,
  );
  expect(screen.getByRole('status')).toHaveTextContent('Loading records');
  expect(screen.getByRole('heading', { name: 'Summary' })).toBeInTheDocument();
  expect(screen.getByText('Version: 3')).toBeInTheDocument();
  expect(screen.getByText('Date: 2026-10-25')).toBeInTheDocument();
});

it('uses the common login hierarchy with pending and safe error feedback', () => {
  const profile = structuredClone(identity);
  profile.application.name = 'Programme Portal';
  const { container } = render(
    <IdentityShell authentication profile={profile}>
      <IdentitySignIn profile={profile} pending error="Please try again" />
    </IdentityShell>,
  );
  expect(
    screen.getByRole('heading', { name: 'Sign in', level: 1 }),
  ).toBeTruthy();
  expect(
    screen.getByRole('button', { name: 'Opening sign-in…' }),
  ).toBeDisabled();
  expect(screen.getByRole('alert')).toHaveTextContent('Please try again');
  expect(container.querySelector('.identity-login')).toHaveAttribute(
    'aria-busy',
    'true',
  );
  expect(container.querySelector('.identity-header')).toBeNull();
  expect(
    screen.getByText('Use your organisation account to continue.'),
  ).toBeTruthy();
});
