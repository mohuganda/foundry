import { createHomePage } from '../pages/home/home-page.js';
import {
  createIdentityShell,
  createIdentityState,
} from '../shared/identity/dom.js';
import { auth } from '../shared/auth/client.js';

export function mountApp(root) {
  let cleanup = () => {};
  const render = () => {
    cleanup();
    const session = auth.getSession();
    const content = document.createElement('div');
    let view;
    if (session.status === 'loading') view = createIdentityState('loading');
    else if (session.status === 'error')
      view = createIdentityState('error', {
        message: session.error,
        onAction: () => auth.signIn(),
      });
    else if (session.status === 'anonymous')
      view = createIdentityState('sign-in', { onAction: () => auth.signIn() });
    else {
      view = createHomePage();
      if (session.status === 'authenticated') {
        const button = document.createElement('button');
        button.className = 'button';
        button.type = 'button';
        button.textContent = 'Sign out';
        button.addEventListener('click', () => {
          void auth.signOut();
        });
        content.append(button);
      }
    }
    content.append(view.element);
    const shell = createIdentityShell(content);
    root.replaceChildren(shell.element);
    cleanup = () => {
      view.destroy();
      shell.destroy();
    };
  };
  const unsubscribe = auth.subscribe(render);
  render();
  void auth.initialize();
  return () => {
    unsubscribe();
    cleanup();
  };
}
