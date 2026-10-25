import { useState } from 'react';
import {
  assetURL,
  identity,
  publicLinks,
  reportMetadata,
  safeHref,
} from './index.js';

function BrandImage({ src, dark, alt }) {
  const [failed, setFailed] = useState(false);
  if (!src || failed) return null;
  return (
    <picture>
      {dark && <source media="(prefers-color-scheme: dark)" srcSet={dark} />}
      <img src={src} alt={alt} onError={() => setFailed(true)} />
    </picture>
  );
}

function BrandImages({ profile }) {
  const assets = profile.branding.assets;
  const logo = assetURL(assets.logo_light || assets.logo_dark);
  const dark = assetURL(assets.logo_dark);
  const arms = assetURL(assets.coat_of_arms);
  return (
    <span className="identity-images">
      <BrandImage
        key={`${logo}:${dark}`}
        src={logo}
        dark={dark}
        alt={`${profile.organisation.name} logo`}
      />
      <BrandImage key={arms} src={arms} alt="Coat of arms" />
    </span>
  );
}

export function IdentityShell({ children, profile = identity }) {
  const context = [
    profile.application.department,
    profile.application.programme,
  ].filter(Boolean);
  const navigation = profile.navigation.filter((link) => safeHref(link.href));
  return (
    <div className="shell identity-shell">
      <a className="skip-link" href="#main">
        Skip to content
      </a>
      <header className="identity-header">
        <BrandImages profile={profile} />
        <div>
          {profile.organisation.country && (
            <p>{profile.organisation.country}</p>
          )}
          <strong>{profile.organisation.name}</strong>
          <p>{profile.application.name}</p>
          {context.length > 0 && <p className="muted">{context.join(' / ')}</p>}
        </div>
        {profile.application.environment !== 'production' && (
          <span className="identity-environment">
            {profile.application.environment} environment
          </span>
        )}
      </header>
      {navigation.length > 0 && (
        <nav aria-label="Main">
          {navigation.map((link) => (
            <a key={`${link.label}:${link.href}`} href={link.href}>
              {link.label}
            </a>
          ))}
        </nav>
      )}
      <main id="main" tabIndex={-1}>
        {children}
      </main>
      <footer className="identity-footer muted">
        <p>
          {profile.organisation.name} · {profile.application.name}
        </p>
        <nav aria-label="Information and support">
          {publicLinks(profile).map((link) => (
            <a key={link.label} href={link.href}>
              {link.label}
            </a>
          ))}
        </nav>
      </footer>
    </div>
  );
}

export function IdentityLoading({ message = 'Loading…' }) {
  return (
    <section className="card identity-state">
      <h1>Loading</h1>
      <p role="status">{message}</p>
    </section>
  );
}

export function IdentityError({ message = 'Please try again.', onRetry }) {
  return (
    <section className="card identity-state">
      <h1>Something went wrong</h1>
      <p role="alert">{message}</p>
      <button
        type="button"
        className="button"
        disabled={!onRetry}
        onClick={onRetry}
      >
        Try again
      </button>
    </section>
  );
}

export function IdentitySignIn({ onSignIn, profile = identity }) {
  return (
    <section className="card identity-state">
      <h1>Sign in</h1>
      <p>{profile.application.name}</p>
      <button
        type="button"
        className="button"
        disabled={!onSignIn}
        onClick={onSignIn}
      >
        Continue
      </button>
    </section>
  );
}

export function IdentityReportHeader({
  title,
  version = '',
  issuedOn = '',
  profile = identity,
}) {
  const metadata = reportMetadata(title, { version, issuedOn }, profile);
  return (
    <header className="identity-report-header">
      <BrandImages profile={profile} />
      {metadata.country && <p>{metadata.country}</p>}
      <p>{metadata.organisation}</p>
      <p>{metadata.application}</p>
      <h1>{metadata.title}</h1>
      <p>Version: {metadata.version}</p>
      {metadata.issuedOn && <p>Date: {metadata.issuedOn}</p>}
      {metadata.footer && <p>{metadata.footer}</p>}
    </header>
  );
}
