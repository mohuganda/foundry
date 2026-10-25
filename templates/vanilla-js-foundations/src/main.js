import { applyDocumentIdentity } from './shared/identity/index.js';
import { mountApp } from './app/app.js';
import './shared/styles/global.css';

applyDocumentIdentity();

const root = document.querySelector('#app');
if (!root) throw new Error('Application root #app is missing.');
const unmount = mountApp(root);

if (import.meta.hot) import.meta.hot.dispose(unmount);
