class AuthSession {
  const AuthSession({required this.accessToken, this.userId});
  final String accessToken;
  final String? userId;
}
