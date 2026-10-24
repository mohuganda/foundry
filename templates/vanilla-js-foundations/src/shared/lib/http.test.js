import { afterEach, describe, expect, it, vi } from 'vitest';
import { HttpError, requestJSON } from './http.js';

afterEach(() => vi.unstubAllGlobals());

describe('requestJSON', () => {
  it('returns JSON and forwards cancellation and custom headers', async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValue(
        new Response(JSON.stringify({ name: 'Ada' }), { status: 200 }),
      );
    vi.stubGlobal('fetch', fetchMock);
    const controller = new AbortController();

    await expect(
      requestJSON('/api/profile', {
        signal: controller.signal,
        headers: { 'X-Request-ID': 'example' },
      }),
    ).resolves.toEqual({ name: 'Ada' });
    const [, options] = fetchMock.mock.calls[0];
    expect(options.signal).toBe(controller.signal);
    expect(options.headers.get('Accept')).toBe('application/json');
    expect(options.headers.get('X-Request-ID')).toBe('example');
  });

  it('accepts an empty successful response', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(new Response(null, { status: 204 })),
    );
    await expect(requestJSON('/api/profile')).resolves.toBeNull();
  });

  it('reports unsuccessful HTTP responses instead of treating them as data', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(new Response(null, { status: 403 })),
    );
    await expect(requestJSON('/api/profile')).rejects.toBeInstanceOf(HttpError);
    await expect(requestJSON('/api/profile')).rejects.toMatchObject({
      status: 403,
    });
  });

  it('preserves network and cancellation errors for the caller', async () => {
    const error = new DOMException('Cancelled', 'AbortError');
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(error));
    await expect(requestJSON('/api/profile')).rejects.toBe(error);
  });
});
