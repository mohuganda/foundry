import { describe, expect, it, vi } from 'vitest';
import { createAuth } from './client.js';

const settings = {
  mode: 'oidc',
  issuer: 'https://auth.example.org/realms/test',
  audience: 'api',
  scopes: ['openid'],
  web: {
    client_id: 'web',
    redirect_uri: 'https://portal.example.org/auth/callback',
    post_logout_redirect_uri: 'https://portal.example.org/auth/signed-out',
  },
};
const user = {
  profile: { sub: 'person', name: 'User' },
  access_token: 'access-secret',
  refresh_token: 'refresh-secret',
  id_token: 'id-secret',
  expires_in: 300,
  expired: false,
};
function fixture(path = '/', overrides = {}) {
  let stored = { ...user };
  const manager = {
    events: {
      addUserLoaded: vi.fn(),
      addUserUnloaded: vi.fn(),
      addAccessTokenExpired: vi.fn(),
    },
    clearStaleState: vi.fn(),
    getUser: vi.fn(async () => stored),
    removeUser: vi.fn(async () => {
      stored = null;
    }),
    signinRedirect: vi.fn(),
    signoutRedirect: vi.fn(),
    signinRedirectCallback: vi.fn(async () => user),
    signoutRedirectCallback: vi.fn(),
    signinSilent: vi.fn(async () => user),
    ...overrides,
  };
  const factory = vi.fn(() => manager);
  const history = { replaceState: vi.fn() };
  const client = createAuth(settings, {
    managerFactory: factory,
    location: new URL(path, 'https://portal.example.org'),
    history,
  });
  return { client, manager, factory, history };
}
describe('OIDC adapter', () => {
  it('does not create a provider or request login when disabled', async () => {
    const factory = vi.fn();
    const client = createAuth({ mode: 'none' }, { managerFactory: factory });
    await client.initialize();
    expect(client.getSession().status).toBe('disabled');
    expect(await client.getAccessToken()).toBeNull();
    expect(factory).not.toHaveBeenCalled();
  });
  it('uses PKCE and memory token storage, without leaking tokens to session UI', async () => {
    const { client, factory } = fixture();
    await Promise.all([client.initialize(), client.initialize()]);
    expect(factory).toHaveBeenCalledOnce();
    expect(factory.mock.calls[0][0]).toMatchObject({
      response_type: 'code',
      disablePKCE: false,
      automaticSilentRenew: false,
    });
    expect(client.getSession().user.subject).toBe('person');
    expect(JSON.stringify(client.getSession())).not.toContain('secret');
    expect(await client.getAccessToken()).toBe('access-secret');
  });
  it('validates callbacks through the library before clearing the URL', async () => {
    const { client, manager, history } = fixture(
      '/auth/callback?code=code&state=state',
    );
    await client.initialize();
    expect(manager.signinRedirectCallback).toHaveBeenCalledOnce();
    expect(history.replaceState).toHaveBeenCalled();
    expect(client.getSession().status).toBe('authenticated');
  });
  it('fails closed on cancelled or invalid callbacks without exposing provider errors', async () => {
    const { client, manager } = fixture('/auth/callback?error=access_denied', {
      signinRedirectCallback: vi.fn(async () => {
        throw new Error('provider-sensitive-detail');
      }),
    });
    await client.initialize();
    expect(manager.removeUser).toHaveBeenCalled();
    expect(client.getSession().status).toBe('error');
    expect(client.getSession().error).not.toContain('sensitive');
  });
  it('deduplicates refresh and clears invalid sessions', async () => {
    const { client, manager } = fixture('/', {
      getUser: vi.fn(async () => ({ ...user, expires_in: 10 })),
      signinSilent: vi.fn(async () => {
        throw new Error('invalid_grant');
      }),
    });
    await client.initialize();
    expect(
      await Promise.all([client.getAccessToken(), client.getAccessToken()]),
    ).toEqual([null, null]);
    expect(manager.signinSilent).toHaveBeenCalledOnce();
    expect(client.getSession().status).toBe('error');
  });
  it('logs out locally before redirect and handles logout callbacks', async () => {
    const { client, manager } = fixture();
    await client.initialize();
    await client.signOut();
    expect(manager.removeUser).toHaveBeenCalled();
    expect(manager.signoutRedirect).toHaveBeenCalledWith({
      id_token_hint: 'id-secret',
    });
    expect(client.getSession().status).toBe('anonymous');
    const logout = fixture('/auth/signed-out?state=state');
    await logout.client.initialize();
    expect(logout.manager.signoutRedirectCallback).toHaveBeenCalled();
    expect(logout.client.getSession().status).toBe('anonymous');
  });
  it('rejects a mismatched application origin', async () => {
    const client = createAuth(settings, {
      location: new URL('https://other.example.org'),
    });
    await client.initialize();
    expect(client.getSession().status).toBe('error');
  });
  it('cannot restore a session when refresh finishes after sign-out', async () => {
    let resolve;
    const pending = new Promise((done) => {
      resolve = done;
    });
    const { client, manager } = fixture('/', {
      getUser: vi.fn(async () => ({ ...user, expires_in: 10 })),
      signinSilent: vi.fn(() => pending),
    });
    await client.initialize();
    const token = client.getAccessToken();
    await vi.waitFor(() => expect(manager.signinSilent).toHaveBeenCalledOnce());
    await client.signOut();
    manager.events.addUserLoaded.mock.calls[0][0](user);
    resolve(user);
    expect(await token).toBeNull();
    expect(client.getSession().status).toBe('anonymous');
    expect(await client.getAccessToken()).toBeNull();
  });
});
