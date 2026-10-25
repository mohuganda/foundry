import { applyDocumentIdentity } from './shared/identity/index.js';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { App } from './app/App.jsx';
import './shared/styles/global.css';

applyDocumentIdentity();

const root = document.querySelector('#app');
if (!root) throw new Error('Application root #app is missing.');
createRoot(root).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
