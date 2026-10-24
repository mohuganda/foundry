import 'package:flutter_secure_storage/flutter_secure_storage.dart';

final class SecureTokenStore {
  const SecureTokenStore(this._storage);
  final FlutterSecureStorage _storage;
  static const accessTokenKey = 'auth.access_token';
  static const refreshTokenKey = 'auth.refresh_token';
  Future<String?> readAccessToken() => _storage.read(key: accessTokenKey);
  Future<String?> readRefreshToken() => _storage.read(key: refreshTokenKey);
  Future<void> writeTokens({
    required String accessToken,
    String? refreshToken,
  }) async {
    await _storage.write(key: accessTokenKey, value: accessToken);
    if (refreshToken != null) {
      await _storage.write(key: refreshTokenKey, value: refreshToken);
    }
  }

  Future<void> clear() async {
    await _storage.delete(key: accessTokenKey);
    await _storage.delete(key: refreshTokenKey);
  }
}
