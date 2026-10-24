import { useState } from 'react';
import { Button } from '../../../shared/ui/Button.jsx';

export function Counter() {
  const [count, setCount] = useState(0);

  return (
    <section className="card" aria-label="Example counter">
      <h2>Your first feature</h2>
      <p>An isolated component with accessible controls and its own tests.</p>
      <output className="counter-value" aria-label="Count" aria-live="polite">
        {count}
      </output>
      <div className="actions">
        <Button onClick={() => setCount((current) => current + 1)}>
          Increase count
        </Button>
        <Button onClick={() => setCount(0)}>Reset</Button>
      </div>
    </section>
  );
}
