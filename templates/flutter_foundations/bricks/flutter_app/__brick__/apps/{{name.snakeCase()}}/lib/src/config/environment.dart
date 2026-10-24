enum AppFlavor { development, staging, production }

final class AppEnvironment {
  const AppEnvironment({
    required this.flavor,
    required this.apiBaseUrl,
    required this.authBaseUrl,
    required this.enableHttpLogs,
  });
  final AppFlavor flavor;
  final String apiBaseUrl;
  final String authBaseUrl;
  final bool enableHttpLogs;

  factory AppEnvironment.fromDartDefines() {
    const env = String.fromEnvironment('APP_ENV', defaultValue: 'development');
    final flavor = switch (env) {
      'production' => AppFlavor.production,
      'staging' => AppFlavor.staging,
      _ => AppFlavor.development,
    };
    return AppEnvironment(
      flavor: flavor,
      apiBaseUrl: const String.fromEnvironment(
        'API_BASE_URL',
        defaultValue: 'http://localhost:9000',
      ),
      authBaseUrl: const String.fromEnvironment(
        'AUTH_BASE_URL',
        defaultValue: 'http://localhost:8080',
      ),
      enableHttpLogs: const bool.fromEnvironment(
        'ENABLE_HTTP_LOGS',
        defaultValue: true,
      ),
    );
  }
}
