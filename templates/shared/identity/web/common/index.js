import config from './config.json';

// Generated public configuration. Edit identity/identity.json and run identity sync.
export const identity = Object.freeze(config);

export function safeHref(value) {
  if (typeof value !== 'string' || /[\\\s]/u.test(value)) return undefined;
  if (/^https:\/\//u.test(value)) {
    try {
      const url = new URL(value);
      return url.username || url.password ? undefined : value;
    } catch {
      return undefined;
    }
  }
  return /^\/(?!\/)|^#/u.test(value) ? value : undefined;
}

export function assetURL(value) {
  if (
    !/^\/branding\/assets\/[a-z_]+\.(svg|png|jpe?g|ico)$/u.test(value || '')
  ) {
    return undefined;
  }
  return `${import.meta.env.BASE_URL}${value.slice(1)}`;
}

export function publicLinks(profile = identity) {
  const links = [
    [
      profile.profile === 'moh-uganda'
        ? 'Ministry website'
        : 'Organisation website',
      profile.organisation.website,
    ],
    ['Privacy', profile.links.privacy],
    ['Accessibility', profile.links.accessibility],
    ['Support', profile.links.support],
  ];
  const result = links
    .filter(([, href]) => safeHref(href))
    .map(([label, href]) => ({ label, href }));
  const email = profile.application.support_email;
  if (email && /^[^\s?@]+@[^\s?@]+\.[^\s?@]+$/u.test(email)) {
    result.push({
      label: 'Email support',
      href: `mailto:${encodeURIComponent(email)}`,
    });
  }
  return result;
}

export function reportMetadata(
  title,
  { version = '', issuedOn = '' } = {},
  profile = identity,
) {
  return {
    organisation: profile.organisation.name,
    country: profile.organisation.country,
    application: profile.application.name,
    title: `${profile.documents.title_prefix}${title}`,
    version: version || profile.documents.default_version,
    issuedOn,
    footer: profile.documents.footer_text,
    exportTitle: `${profile.application.name} - ${title}`,
  };
}

export function applyDocumentIdentity(doc = document, profile = identity) {
  doc.title = profile.application.name;
  const description = doc.querySelector('meta[name="description"]');
  if (description) description.content = profile.application.description;
  for (const [rel, value] of [
    ['icon', profile.branding.assets.favicon],
    ['apple-touch-icon', profile.branding.assets.application_icon],
  ]) {
    doc.querySelector(`link[data-foundry-identity="${rel}"]`)?.remove();
    const href = assetURL(value);
    if (!href) continue;
    const link = doc.createElement('link');
    link.rel = rel;
    link.href = href;
    link.dataset.foundryIdentity = rel;
    doc.head.append(link);
  }
}
