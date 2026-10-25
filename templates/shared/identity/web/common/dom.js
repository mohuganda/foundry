import {
  assetURL,
  identity,
  publicLinks,
  reportMetadata,
  safeHref,
} from './index.js';

function text(tag, value, className = '') {
  const node = document.createElement(tag);
  node.textContent = value;
  if (className) node.className = className;
  return node;
}

function brandImages(profile) {
  const group = document.createElement('span');
  group.className = 'identity-images';
  const assets = profile.branding.assets;
  const logo = assetURL(assets.logo_light || assets.logo_dark);
  if (logo) {
    const picture = document.createElement('picture');
    const dark = assetURL(assets.logo_dark);
    if (dark) {
      const source = document.createElement('source');
      source.media = '(prefers-color-scheme: dark)';
      source.srcset = dark;
      picture.append(source);
    }
    const img = document.createElement('img');
    img.src = logo;
    img.alt = `${profile.organisation.name} logo`;
    img.addEventListener('error', () => picture.remove(), { once: true });
    picture.append(img);
    group.append(picture);
  }
  const arms = assetURL(assets.coat_of_arms);
  if (arms) {
    const img = document.createElement('img');
    img.src = arms;
    img.alt = 'Coat of arms';
    img.addEventListener('error', () => img.remove(), { once: true });
    group.append(img);
  }
  return group;
}

export function createIdentityShell(content, profile = identity) {
  const element = document.createElement('div');
  element.className = 'shell identity-shell';
  const skip = text('a', 'Skip to content', 'skip-link');
  skip.href = '#main';
  const header = document.createElement('header');
  header.className = 'identity-header';
  header.append(brandImages(profile));
  const label = document.createElement('div');
  if (profile.organisation.country)
    label.append(text('p', profile.organisation.country));
  label.append(text('strong', profile.organisation.name));
  label.append(text('p', profile.application.name));
  const context = [
    profile.application.department,
    profile.application.programme,
  ].filter(Boolean);
  if (context.length) label.append(text('p', context.join(' / '), 'muted'));
  header.append(label);
  if (profile.application.environment !== 'production') {
    header.append(
      text(
        'span',
        `${profile.application.environment} environment`,
        'identity-environment',
      ),
    );
  }
  const nav = document.createElement('nav');
  nav.setAttribute('aria-label', 'Main');
  for (const link of profile.navigation) {
    const href = safeHref(link.href);
    if (!href) continue;
    const anchor = text('a', link.label);
    anchor.href = href;
    nav.append(anchor);
  }
  const main = document.createElement('main');
  main.id = 'main';
  main.tabIndex = -1;
  main.append(content);
  const footer = document.createElement('footer');
  footer.className = 'identity-footer muted';
  footer.append(
    text('p', `${profile.organisation.name} · ${profile.application.name}`),
  );
  const links = document.createElement('nav');
  links.setAttribute('aria-label', 'Information and support');
  for (const link of publicLinks(profile)) {
    const anchor = text('a', link.label);
    anchor.href = link.href;
    links.append(anchor);
  }
  footer.append(links);
  element.append(skip, header);
  if (nav.childElementCount) element.append(nav);
  element.append(main, footer);
  return { element, destroy: () => element.remove() };
}

export function createIdentityState(
  kind,
  { message = '', onAction } = {},
  profile = identity,
) {
  const element = document.createElement('section');
  element.className = 'card identity-state';
  const heading = {
    loading: 'Loading',
    error: 'Something went wrong',
    'sign-in': 'Sign in',
  }[kind];
  if (!heading) throw new Error('Unsupported identity state.');
  element.append(text('h1', heading));
  const status = text(
    'p',
    message || (kind === 'sign-in' ? profile.application.name : heading),
  );
  status.setAttribute('role', kind === 'error' ? 'alert' : 'status');
  element.append(status);
  const controller = new AbortController();
  if (kind !== 'loading') {
    const button = text(
      'button',
      kind === 'sign-in' ? 'Continue' : 'Try again',
      'button',
    );
    button.type = 'button';
    button.disabled = !onAction;
    if (onAction)
      button.addEventListener('click', onAction, { signal: controller.signal });
    element.append(button);
  }
  return {
    element,
    destroy() {
      controller.abort();
      element.remove();
    },
  };
}

export function createIdentityReportHeader(
  title,
  options = {},
  profile = identity,
) {
  const metadata = reportMetadata(title, options, profile);
  const element = document.createElement('header');
  element.className = 'identity-report-header';
  element.append(brandImages(profile));
  for (const value of [
    metadata.country,
    metadata.organisation,
    metadata.application,
  ]) {
    if (value) element.append(text('p', value));
  }
  element.append(text('h1', metadata.title));
  element.append(text('p', `Version: ${metadata.version}`));
  if (metadata.issuedOn)
    element.append(text('p', `Date: ${metadata.issuedOn}`));
  if (metadata.footer) element.append(text('p', metadata.footer));
  return element;
}
