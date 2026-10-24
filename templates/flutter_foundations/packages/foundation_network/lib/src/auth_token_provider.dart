abstract interface class AuthTokenProvider {
  Future<String?> accessToken();
  Future<String?> refreshAccessToken();
  Future<void> clearSession();
}
