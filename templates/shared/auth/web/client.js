import {
  InMemoryWebStorage,
  UserManager,
  WebStorageStateStore,
} from 'oidc-client-ts';
import config from './config.json';

// Tokens stay in memory. Only redirect transaction state uses sessionStorage.
export function createAuth(settings, dependencies = {}) {
  const location = dependencies.location ?? window.location;
  const history = dependencies.history ?? window.history;
  let manager;
  let initializing;
  let refreshing;
  let generation = 0;
  let acceptingSession = true;
  let snapshot = Object.freeze({
    status: settings.mode === 'none' ? 'disabled' : 'loading',
    user: null,
    error: '',
  });
  const listeners = new Set();
  const publish = (status, user = null, error = '') => {
    snapshot = Object.freeze({
      status,
      user: user
        ? Object.freeze({
            subject: user.profile.sub,
            name: user.profile.name ?? user.profile.preferred_username ?? '',
            email: user.profile.email ?? '',
          })
        : null,
      error,
    });
    for (const listener of listeners) listener();
  };
  const failed = async () => {
    acceptingSession = false;
    try {
      await manager?.removeUser();
    } catch {
      /* Session remains unavailable. */
    }
    publish(
      'error',
      null,
      'Authentication could not be completed. Please sign in again.',
    );
  };
  const getManager = () => {
    if (manager) return manager;
    if (settings.mode !== 'oidc') throw new Error('OIDC is not enabled.');
    const redirect = new URL(settings.web.redirect_uri);
    const logout = new URL(settings.web.post_logout_redirect_uri);
    if (
      redirect.origin !== location.origin ||
      logout.origin !== location.origin
    )
      throw new Error(
        'Authentication callback origin does not match this application.',
      );
    const options = {
      authority: settings.issuer,
      client_id: settings.web.client_id,
      redirect_uri: redirect.href,
      post_logout_redirect_uri: logout.href,
      response_type: 'code',
      response_mode: 'query',
      scope: settings.scopes.join(' '),
      disablePKCE: false,
      automaticSilentRenew: false,
      monitorSession: false,
      loadUserInfo: false,
      userStore: new WebStorageStateStore({ store: new InMemoryWebStorage() }),
      stateStore: new WebStorageStateStore({
        store: dependencies.stateStorage ?? window.sessionStorage,
      }),
    };
    manager = dependencies.managerFactory
      ? dependencies.managerFactory(options)
      : new UserManager(options);
    manager.events.addUserLoaded((user) => {
      if (acceptingSession) publish('authenticated', user);
    });
    manager.events.addUserUnloaded(() => publish('anonymous'));
    manager.events.addAccessTokenExpired(() => {
      void failed();
    });
    return manager;
  };
  return {
    getSession: () => snapshot,
    subscribe(listener) {
      listeners.add(listener);
      return () => listeners.delete(listener);
    },
    initialize() {
      if (initializing) return initializing;
      initializing = (async () => {
        if (settings.mode === 'none') return;
        const started = generation;
        try {
          const mgr = getManager();
          await mgr.clearStaleState();
          const pathname = location.pathname;
          let user;
          if (pathname === new URL(settings.web.redirect_uri).pathname) {
            user = await mgr.signinRedirectCallback(location.href);
            history.replaceState({}, '', import.meta.env.BASE_URL);
          } else if (
            pathname === new URL(settings.web.post_logout_redirect_uri).pathname
          ) {
            await mgr.signoutRedirectCallback(location.href);
            await mgr.removeUser();
            history.replaceState({}, '', import.meta.env.BASE_URL);
          } else user = await mgr.getUser();
          if (started !== generation) {
            await mgr.removeUser();
            return;
          }
          if (user && !user.expired) publish('authenticated', user);
          else {
            await mgr.removeUser();
            publish('anonymous');
          }
        } catch {
          if (started === generation) await failed();
        } finally {
          if (
            [
              settings.web.redirect_uri,
              settings.web.post_logout_redirect_uri,
            ].some((value) => new URL(value).pathname === location.pathname)
          )
            history.replaceState({}, '', import.meta.env.BASE_URL);
        }
      })();
      return initializing;
    },
    async signIn() {
      generation++;
      acceptingSession = true;
      try {
        publish('loading');
        await getManager().signinRedirect();
      } catch {
        await failed();
      }
    },
    async signOut() {
      generation++;
      acceptingSession = false;
      try {
        const mgr = getManager();
        const user = await mgr.getUser();
        await mgr.removeUser();
        publish('anonymous');
        await mgr.signoutRedirect({ id_token_hint: user?.id_token });
      } catch {
        await failed();
      }
    },
    async getAccessToken() {
      if (settings.mode === 'none') return null;
      await this.initialize();
      const started = generation;
      const mgr = getManager();
      const user = await mgr.getUser();
      if (!acceptingSession || started !== generation || !user) return null;
      if (user.expires_in > 30) return user.access_token;
      if (!refreshing)
        refreshing = (async () => {
          try {
            // Do not rely on hidden iframes/third-party cookies as a refresh fallback.
            if (!user.refresh_token) throw new Error('Sign-in required.');
            const renewed = await mgr.signinSilent();
            if (started !== generation || !acceptingSession) {
              await mgr.removeUser();
              return null;
            }
            if (!renewed || renewed.expired)
              throw new Error('Sign-in required.');
            publish('authenticated', renewed);
            return renewed.access_token;
          } catch {
            if (started === generation) await failed();
            return null;
          } finally {
            refreshing = undefined;
          }
        })();
      return refreshing;
    },
  };
}
export const auth = createAuth(config);
