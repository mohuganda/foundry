import { useEffect, useSyncExternalStore } from 'react';
import { HomePage } from '../pages/home/HomePage.jsx';
import { IdentityShell, IdentitySignIn } from '../shared/identity/Identity.jsx';
import { auth } from '../shared/auth/client.js';

export function App() {
  const session = useSyncExternalStore(auth.subscribe, auth.getSession);
  useEffect(() => {
    void auth.initialize();
  }, []);
  let content;
  if (session.status === 'loading') content = <IdentitySignIn pending />;
  else if (session.status === 'error')
    content = (
      <IdentitySignIn error={session.error} onSignIn={() => auth.signIn()} />
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
  return (
    <IdentityShell
      authentication={!['disabled', 'authenticated'].includes(session.status)}
    >
      {content}
    </IdentityShell>
  );
}
