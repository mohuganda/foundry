export function Button({ type = 'button', className = '', ...props }) {
  return (
    <button type={type} className={`button ${className}`.trim()} {...props} />
  );
}
