import { auth } from './client.js';
import { env } from '../config/env.js';
import { requestJSON } from '../lib/http.js';

// Restrict bearer credentials to the configured API origin and base path.
export async function authenticatedRequest(path, { query, ...options } = {}) {
  if (
    typeof path !== 'string' ||
    !path.startsWith('/') ||
    path.startsWith('//') ||
    /[\\?#]/u.test(path) ||
    path
      .split('/')
      .some((part) => part === '.' || part === '..' || part.includes('%'))
  )
    throw new Error('Use an API-relative path without query or traversal.');
  const base = new URL(
    env.apiBaseUrl.endsWith('/') ? env.apiBaseUrl : `${env.apiBaseUrl}/`,
    window.location.origin,
  );
  const target = new URL(path.slice(1), base);
  if (
    target.origin !== base.origin ||
    !target.pathname.startsWith(base.pathname)
  )
    throw new Error('API URL escaped configured base.');
  if (
    target.protocol !== 'https:' &&
    target.hostname !== 'localhost' &&
    target.hostname !== '127.0.0.1'
  )
    throw new Error('API requires HTTPS.');
  if (query) target.search = new URLSearchParams(query).toString();
  const headers = new Headers(options.headers);
  const token = await auth.getAccessToken();
  if (!token) throw new Error('Sign-in required.');
  headers.set('Authorization', `Bearer ${token}`);
  // Redirects must not forward credentials to a different endpoint.
  return requestJSON(target, { ...options, headers, redirect: 'error' });
}
