import { HomePage } from '../pages/home/HomePage.jsx';
import { IdentityShell } from '../shared/identity/Identity.jsx';

export function App() {
  return (
    <IdentityShell>
      <HomePage />
    </IdentityShell>
  );
}
