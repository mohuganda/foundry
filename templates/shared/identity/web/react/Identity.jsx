import { useState } from 'react';
import { CircleAlert, LoaderCircle, LogIn, RotateCcw } from 'lucide-react';
import login from './login.json';
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

export function IdentityShell({
  children,
  profile = identity,
  authentication = false,
}) {
  const context = [
    profile.application.department,
    profile.application.programme,
  ].filter(Boolean);
  const navigation = profile.navigation.filter((link) => safeHref(link.href));
  return (
    <div
      className={`shell identity-shell${authentication ? ' identity-shell--auth' : ''}`}
    >
      <a className="skip-link" href="#main">
        Skip to content
      </a>
      {!authentication && (
        <header className="identity-header">
          <BrandImages profile={profile} />
          <div>
            {profile.organisation.country && (
              <p>{profile.organisation.country}</p>
            )}
            <strong>{profile.organisation.name}</strong>
            <p>{profile.application.name}</p>
            {context.length > 0 && (
              <p className="muted">{context.join(' / ')}</p>
            )}
          </div>
          {profile.application.environment !== 'production' && (
            <span className="identity-environment">
              {profile.application.environment} environment
            </span>
          )}
        </header>
      )}
      {!authentication && navigation.length > 0 && (
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
      <LoaderCircle className="identity-icon" aria-hidden="true" />
      <h1>Loading</h1>
      <p role="status">{message}</p>
    </section>
  );
}

export function IdentityError({ message = 'Please try again.', onRetry }) {
  return (
    <section className="card identity-state">
      <CircleAlert className="identity-icon" aria-hidden="true" />
      <h1>Something went wrong</h1>
      <p role="alert">{message}</p>
      <button
        type="button"
        className="button"
        disabled={!onRetry}
        onClick={onRetry}
      >
        <RotateCcw className="identity-icon" aria-hidden="true" />
        Try again
      </button>
    </section>
  );
}

export function IdentitySignIn({
  onSignIn,
  profile = identity,
  pending = false,
  error = '',
}) {
  const context = [
    profile.application.department,
    profile.application.programme,
  ]
    .filter(Boolean)
    .join(' / ');
  return (
    <section
      className="identity-login"
      aria-labelledby="sign-in-title"
      aria-busy={pending}
      style={{
        '--login-max-width': `${login.max_width}px`,
        '--login-padding': `${login.padding}px`,
        '--login-control-height': `${login.control_height}px`,
        '--login-title-size': `${login.title_size}px`,
      }}
    >
      <div className="identity-login-brand">
        <BrandImages profile={profile} />
        {profile.organisation.country && <p>{profile.organisation.country}</p>}
        <strong>{profile.organisation.name}</strong>
        {profile.application.environment !== 'production' && (
          <span className="identity-environment">
            {profile.application.environment} environment
          </span>
        )}
      </div>
      <div className="identity-login-content">
        <p className="identity-login-eyebrow">{login.eyebrow}</p>
        <p className="identity-login-product">{profile.application.name}</p>
        {context && <p className="muted">{context}</p>}
        <h1 id="sign-in-title">{login.title}</h1>
        <p className="muted">{login.description}</p>
        {error && (
          <p className="identity-login-error" role="alert">
            <CircleAlert className="identity-icon" aria-hidden="true" />
            {error}
          </p>
        )}
        <button
          type="button"
          className="button identity-login-action"
          disabled={pending || !onSignIn}
          onClick={onSignIn}
        >
          {pending ? (
            <LoaderCircle className="identity-icon" aria-hidden="true" />
          ) : (
            <LogIn className="identity-icon" aria-hidden="true" />
          )}
          {pending ? login.pending_action : login.action}
        </button>
        {pending && (
          <p className="identity-login-progress" role="status">
            {login.pending_action}
          </p>
        )}
        <p className="identity-login-notice muted">{login.notice}</p>
      </div>
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
