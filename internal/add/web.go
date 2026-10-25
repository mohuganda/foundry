package add

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jabahum/foundry/internal/workspace"
)

// WebFeature creates a UI boundary and colocated tests without editing pages.
func WebFeature(root, name string, options ...Options) (Result, error) {
	name, err := featureName(name)
	if err != nil {
		return Result{}, err
	}
	p, err := workspace.At(root)
	if err != nil {
		return Result{}, err
	}
	if p.Kind != "javascript" && p.Kind != "react" {
		return Result{}, fmt.Errorf("web feature generation requires a JavaScript or React workspace")
	}
	if len(options) > 0 && options[0].App != "" {
		return Result{}, fmt.Errorf("--app is only supported for Flutter features")
	}
	title, component := featureTitle(name), pascalCase(name)
	base := filepath.Join("src", "features", name)
	render := strings.NewReplacer("__TITLE__", title, "__COMPONENT__", component, "__NAME__", name)
	files := map[string]string{}
	if p.Kind == "react" {
		files[filepath.Join(base, "components", component+".jsx")] = render.Replace(`export function __COMPONENT__({ message = 'Ready to build.' }) {
  return (
    <section className="card" aria-label="__TITLE__">
      <h2>__TITLE__</h2>
      <p role="status">{message}</p>
    </section>
  );
}
`)
		files[filepath.Join(base, "components", component+".test.jsx")] = render.Replace(`import { expect, it } from 'vitest';
import { render, within } from '@testing-library/react';
import { __COMPONENT__ } from './__COMPONENT__.jsx';

it('updates the visible message when the page supplies new state', () => {
  const { getByRole, rerender } = render(<__COMPONENT__ message="Loading" />);
  const feature = getByRole('region', { name: '__TITLE__' });
  expect(within(feature).getByRole('status')).toHaveTextContent('Loading');

  rerender(<__COMPONENT__ message="Loaded" />);
  expect(within(feature).getByRole('status')).toHaveTextContent('Loaded');
});

it('keeps separate feature instances independent', () => {
  const { getAllByRole } = render(
    <>
      <__COMPONENT__ message="First" />
      <__COMPONENT__ message="Second" />
    </>,
  );
  const [first, second] = getAllByRole('region', { name: '__TITLE__' });
  expect(within(first).getByRole('status')).toHaveTextContent('First');
  expect(within(second).getByRole('status')).toHaveTextContent('Second');
});
`)
	} else {
		files[filepath.Join(base, name+".js")] = render.Replace(`// Compose this feature from a page; keep product behavior inside this folder.
export function create__COMPONENT__({ message = 'Ready to build.' } = {}) {
  const element = document.createElement('section');
  element.className = 'card';
  element.setAttribute('aria-label', '__TITLE__');
  const heading = document.createElement('h2');
  heading.textContent = '__TITLE__';
  const status = document.createElement('p');
  status.setAttribute('role', 'status');
  status.textContent = message;
  element.append(heading, status);

  return {
    element,
    update(message) {
      status.textContent = message;
    },
    destroy() {
      element.remove();
    },
  };
}
`)
		files[filepath.Join(base, name+".test.js")] = render.Replace(`import { expect, it } from 'vitest';
import { getByRole } from '@testing-library/dom';
import { create__COMPONENT__ } from './__NAME__.js';

it('updates visible state and removes its DOM on teardown', () => {
  const feature = create__COMPONENT__({ message: 'Loading' });
  document.body.append(feature.element);
  expect(getByRole(feature.element, 'status')).toHaveTextContent('Loading');

  feature.update('Loaded');
  expect(getByRole(feature.element, 'status')).toHaveTextContent('Loaded');
  feature.destroy();
  expect(feature.element).not.toBeInTheDocument();
});

it('keeps separate feature instances independent', () => {
  const first = create__COMPONENT__({ message: 'First' });
  const second = create__COMPONENT__({ message: 'Second' });
  document.body.append(first.element, second.element);
  first.update('Changed');
  expect(getByRole(second.element, 'status')).toHaveTextContent('Second');
  first.destroy();
  second.destroy();
});
`)
	}
	files[filepath.Join(base, "README.md")] = render.Replace(`# __TITLE__

This feature owns its UI and colocated tests. Replace the starter message boundary
with product behavior, keeping reusable code in src/shared/ and composing the
feature from src/pages/. Add feature-local API, state, and styles when needed.

Run npm run check from the project root. The generator does not install packages,
add routes, or modify existing pages. Use the import printed by Foundry to wire it.
`)
	written, err := writeFiles(root, files, p.Kind, options...)
	if err != nil {
		return Result{}, err
	}
	integration := "import { create" + component + " } from '../../features/" + name + "/" + name + ".js'; then mount element and call destroy on teardown"
	if p.Kind == "react" {
		integration = "import { " + component + " } from '../../features/" + name + "/components/" + component + ".jsx'; then render the component"
	}
	return Result{Files: written, Next: []string{"from a page in src/pages/<page>/: " + integration, "replace the starter message boundary with product behavior and tests", "run: npm run check"}}, nil
}
