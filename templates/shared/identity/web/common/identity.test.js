import { expect, it, vi } from 'vitest';
import { getByRole, queryByRole } from '@testing-library/dom';
import userEvent from '@testing-library/user-event';
import {
  applyDocumentIdentity,
  identity,
  reportMetadata,
  safeHref,
} from './index.js';
import {
  createIdentityReportHeader,
  createIdentityShell,
  createIdentityState,
} from './dom.js';

function profile() {
  const value = structuredClone(identity);
  value.organisation.name = 'Ministry of Health';
  value.organisation.country = 'Republic of Uganda';
  value.application.name = 'Programme Portal';
  value.navigation = [{ label: 'Home', href: '/' }];
  value.links.privacy = '/privacy';
  return value;
}

it('renders configured identity and public links with a usable main landmark', () => {
  const value = profile();
  value.ownership = { technical_owner: 'private-reference' };
  const shell = createIdentityShell(document.createElement('section'), value);
  document.body.append(shell.element);
  expect(getByRole(shell.element, 'main')).toHaveAttribute('id', 'main');
  expect(getByRole(shell.element, 'link', { name: 'Privacy' })).toHaveAttribute(
    'href',
    '/privacy',
  );
  expect(shell.element).toHaveTextContent('Republic of Uganda');
  expect(shell.element).not.toHaveTextContent('private-reference');
  shell.destroy();
});

it('treats identity text as text and filters unsafe links', () => {
  const value = profile();
  value.organisation.name = '<script>unsafe()</script>';
  value.navigation = [{ label: 'Bad', href: 'javascript:unsafe()' }];
  const shell = createIdentityShell(document.createElement('div'), value);
  expect(shell.element.querySelector('script')).toBeNull();
  expect(queryByRole(shell.element, 'link', { name: 'Bad' })).toBeNull();
  expect(safeHref('//untrusted.example')).toBeUndefined();
  shell.destroy();
});

it('uses real auth/retry callbacks and tears listeners down', async () => {
  const callback = vi.fn();
  const view = createIdentityState(
    'sign-in',
    { onAction: callback },
    profile(),
  );
  document.body.append(view.element);
  const button = getByRole(view.element, 'button', { name: 'Continue' });
  await userEvent.setup().click(button);
  expect(callback).toHaveBeenCalledTimes(1);
  view.destroy();
  button.click();
  expect(callback).toHaveBeenCalledTimes(1);
});

it('shows accessible loading/error states and explicit report versions/dates', () => {
  const loading = createIdentityState('loading', {
    message: 'Loading records',
  });
  expect(getByRole(loading.element, 'status')).toHaveTextContent(
    'Loading records',
  );
  const error = createIdentityState('error', { message: 'Unavailable' });
  expect(getByRole(error.element, 'alert')).toHaveTextContent('Unavailable');
  expect(getByRole(error.element, 'button')).toBeDisabled();
  const report = createIdentityReportHeader(
    'Summary',
    { version: '2', issuedOn: '2026-10-25' },
    profile(),
  );
  expect(report).toHaveTextContent('Version: 2');
  expect(report).toHaveTextContent('Date: 2026-10-25');
  expect(reportMetadata('Summary', {}, profile()).exportTitle).toBe(
    'Programme Portal - Summary',
  );
});

it('updates browser metadata and omits missing image assets', () => {
  const value = profile();
  applyDocumentIdentity(document, value);
  expect(document.title).toBe('Programme Portal');
  expect(
    document.querySelector('link[data-foundry-identity="icon"]'),
  ).toBeNull();
  value.branding.assets.favicon = '/branding/assets/favicon.svg';
  applyDocumentIdentity(document, value);
  expect(
    document.querySelector('link[data-foundry-identity="icon"]'),
  ).toHaveAttribute('href', '/branding/assets/favicon.svg');
});
