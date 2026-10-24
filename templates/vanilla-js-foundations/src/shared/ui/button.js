export function createButton(label) {
  const button = document.createElement('button');
  button.type = 'button';
  button.className = 'button';
  button.textContent = label;
  return button;
}
