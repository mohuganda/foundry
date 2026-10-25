import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:foundation_auth/foundation_auth.dart';
import 'package:foundation_network/foundation_network.dart';
import 'package:{{name.snakeCase()}}/src/config/environment.dart';

final environmentProvider = Provider<AppEnvironment>(
  (ref) => throw UnimplementedError('Override in bootstrap'),
);
final authServiceProvider = Provider<AuthService>(
  (ref) => throw UnimplementedError('Override in bootstrap'),
);
final apiClientProvider = Provider<ApiClient>(
  (ref) => throw UnimplementedError('Override in bootstrap'),
);

final oidcAuthProvider = Provider<OidcAuth>(
  (ref) => throw UnimplementedError('Override in bootstrap'),
);
