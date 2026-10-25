import { beforeEach, expect, it, vi } from 'vitest';
vi.mock('./client.js', () => ({
  auth: { getAccessToken: vi.fn(async () => 'token') },
}));
vi.mock('../config/env.js', () => ({
  env: { apiBaseUrl: 'https://api.example.org/v1' },
}));
import { auth } from './client.js';
import { authenticatedRequest } from './request.js';
beforeEach(() => {
  vi.restoreAllMocks();
  auth.getAccessToken.mockResolvedValue('token');
});
it('attaches the access token only to configured API requests', async () => {
  const fetch = vi
    .spyOn(globalThis, 'fetch')
    .mockResolvedValue(new Response('{"ok":true}', { status: 200 }));
  await authenticatedRequest('/reports');
  expect(fetch.mock.calls[0][0].href).toBe(
    'https://api.example.org/v1/reports',
  );
  expect(fetch.mock.calls[0][1].headers.get('Authorization')).toBe(
    'Bearer token',
  );
  expect(fetch.mock.calls[0][1].redirect).toBe('error');
});
it('rejects cross-origin, traversal and missing sessions before requesting', async () => {
  const fetch = vi.spyOn(globalThis, 'fetch');
  for (const path of [
    '//other.example/data',
    '/../../data',
    '/%2e%2e/data',
    '/\\other.example',
  ])
    await expect(authenticatedRequest(path)).rejects.toThrow();
  auth.getAccessToken.mockResolvedValue(null);
  await expect(authenticatedRequest('/reports')).rejects.toThrow(
    'Sign-in required',
  );
  expect(fetch).not.toHaveBeenCalled();
});
