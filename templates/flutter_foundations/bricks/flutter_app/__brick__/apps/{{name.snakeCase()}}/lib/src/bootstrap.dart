import 'package:dio/dio.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:foundation_auth/foundation_auth.dart';
import 'package:foundation_network/foundation_network.dart';
import 'package:foundation_storage/foundation_storage.dart';
import 'package:{{name.snakeCase()}}/src/app.dart';
import 'package:{{name.snakeCase()}}/src/config/environment.dart';
import 'package:{{name.snakeCase()}}/src/providers/foundation_providers.dart';

Future<Widget> bootstrap() async {
  final env = AppEnvironment.fromDartDefines();
  const secureStorage = FlutterSecureStorage();
  final tokenStore = SecureTokenStore(secureStorage);
  final authDio = Dio(BaseOptions(baseUrl: env.authBaseUrl));
  final authService = AuthService(authDio: authDio, tokenStore: tokenStore);
  final apiClient = ApiClient(
    baseUrl: env.apiBaseUrl,
    tokenProvider: authService,
    enableLogs: env.enableHttpLogs,
  );
  return ProviderScope(
    overrides: [
      environmentProvider.overrideWithValue(env),
      authServiceProvider.overrideWithValue(authService),
      apiClientProvider.overrideWithValue(apiClient),
    ],
    child: const AppRoot(),
  );
}
