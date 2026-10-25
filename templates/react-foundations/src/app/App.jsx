import { useEffect, useSyncExternalStore } from 'react';
import { HomePage } from '../pages/home/HomePage.jsx';
import {
  IdentityShell,
  IdentityLoading,
  IdentityError,
  IdentitySignIn,
} from '../shared/identity/Identity.jsx';
import { auth } from '../shared/auth/client.js';

export function App() {
  const session = useSyncExternalStore(auth.subscribe, auth.getSession);
  useEffect(() => {
    void auth.initialize();
  }, []);
  let content;
  if (session.status === 'loading') content = <IdentityLoading />;
  else if (session.status === 'error')
    content = (
      <IdentityError message={session.error} onRetry={() => auth.signIn()} />
    );
  else if (session.status === 'anonymous')
    content = <IdentitySignIn onSignIn={() => auth.signIn()} />;
  else
    content = (
      <>
        {session.status === 'authenticated' && (
          <button
            className="button"
            type="button"
            onClick={() => auth.signOut()}
          >
            Sign out
          </button>
        )}
        <HomePage />
      </>
    );
  return <IdentityShell>{content}</IdentityShell>;
}
