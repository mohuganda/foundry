// VITE_ values are public build-time configuration. Never put secrets here.
export const env = Object.freeze({
  appName: '__PROJECT_NAME__',
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL || '/api',
});
