final class AuthSession {
  const AuthSession({
    required this.accessToken,
    this.refreshToken,
    this.expiresAt,
  });
  final String accessToken;
  final String? refreshToken;
  final DateTime? expiresAt;
}
