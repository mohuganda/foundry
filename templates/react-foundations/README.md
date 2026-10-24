# `__PROJECT_NAME__`

React Foundations: a client-side JavaScript application built with Vite, native ES modules, and a feature-first structure.

## Start

Use Node.js 22.18+ on the 22.x line, Node 24, or a newer supported even-numbered release. `.nvmrc` selects Node 22.

```sh
nvm install
npm ci
npm run dev
```

Run `npm run check` before committing. It checks formatting, lint, tests, and the production build. `npm run test:watch` keeps tests running during development. `npm run preview` serves the built `dist/` folder locally; use a static web host to deploy that folder.

## Structure

```text
public/             Static files copied without bundling
src/
  app/              Application startup and composition
  pages/home/       Screen-level composition
  features/counter/ Feature UI, state, and colocated behavior tests
  shared/
    config/         Public environment configuration
    lib/            Framework-independent utilities and tests
    styles/         Design tokens and global styles
    ui/             Reusable, accessible UI primitives
  test/             Shared test setup
```

Keep business logic and state inside the feature that owns them. Pages compose features; features may import shared code; shared code must not import pages or features. Extract shared code when multiple features need it. Keep tests beside the code they exercise. Add feature-local styles and assets beside their consumers instead of putting every file in global folders.

The counter is an example feature: replace it with your product's first feature. Tests exercise user-visible interactions and independent state. Shared HTTP tests cover unsuccessful responses, empty responses, cancellation, and custom headers.

## Configuration and network requests

Copy `.env.example` to `.env.local` when overriding defaults. Read configuration through `src/shared/config/env.js`. Vite exposes `VITE_` values to every browser; keep credentials and server secrets on your backend. Environment values are applied at build time.

Use `requestJSON` from `src/shared/lib/http.js` for JSON endpoints, passing an `AbortSignal` when requests should stop on navigation. Catch `HttpError` at the feature boundary to show an accessible error state. Add feature-specific API modules inside their feature. Do not start network requests during rendering.

## Conventions

- Native buttons, labelled controls, a skip link, and visible keyboard focus are included.
- Keep state local; introduce shared state and routing when the application requires them.
- Use functional state updates, keep rendering pure, and clean up external effects. React StrictMode and Hooks linting are enabled. Routing, server rendering, and data caching are deliberate additions when needed.
- Dependencies and `package-lock.json` are versioned together. Use `npm ci` for reproducible installs and review dependency updates before regenerating the lockfile.
- `src/main.jsx` is the entry point; `src/app/` owns application composition.

## Foundry delivery

This template is embedded directly in Foundry. Generation needs no network connection; installing dependencies requires npm registry access or a populated npm cache. No dependencies are installed automatically.
