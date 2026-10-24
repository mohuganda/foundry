import 'package:dio/dio.dart';
import 'package:foundation_auth/src/auth_session.dart';
import 'package:foundation_network/foundation_network.dart';
import 'package:foundation_storage/foundation_storage.dart';

final class AuthService implements AuthTokenProvider {
  AuthService({
    required Dio authDio,
    required SecureTokenStore tokenStore,
    this.refreshPath = '/auth/refresh',
  }) : _dio = authDio,
       _tokens = tokenStore;
  final Dio _dio;
  final SecureTokenStore _tokens;
  final String refreshPath;

  @override
  Future<String?> accessToken() => _tokens.readAccessToken();
  @override
  Future<void> clearSession() => _tokens.clear();

  Future<void> persist(AuthSession session) => _tokens.writeTokens(
    accessToken: session.accessToken,
    refreshToken: session.refreshToken,
  );

  @override
  Future<String?> refreshAccessToken() async {
    final refresh = await _tokens.readRefreshToken();
    if (refresh == null || refresh.isEmpty) return null;
    final response = await _dio.post<Map<String, dynamic>>(
      refreshPath,
      data: {'refresh_token': refresh},
    );
    final data = response.data;
    final access = data?['access_token'] as String?;
    if (access == null || access.isEmpty) return null;
    await _tokens.writeTokens(
      accessToken: access,
      refreshToken: data?['refresh_token'] as String? ?? refresh,
    );
    return access;
  }
}
