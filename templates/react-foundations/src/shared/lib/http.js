export class HttpError extends Error {
  constructor(response) {
    super(`Request failed (${response.status})`);
    this.name = 'HttpError';
    this.status = response.status;
  }
}

// Pass an AbortSignal through options to cancel work when a view unmounts.
export async function requestJSON(url, options = {}) {
  const headers = new Headers(options.headers);
  if (!headers.has('Accept')) headers.set('Accept', 'application/json');

  const response = await fetch(url, { ...options, headers });
  if (!response.ok) throw new HttpError(response);
  if (response.status === 204) return null;
  return response.json();
}
