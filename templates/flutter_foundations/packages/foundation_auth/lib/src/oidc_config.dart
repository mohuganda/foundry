import 'dart:convert';
import 'package:foundation_auth/src/generated_config.dart';

/// Non-secret public-client settings. Native redirect registration is required.
class OidcConfig {
  const OidcConfig({
    this.mode = 'none',
    this.issuer = '',
    this.clientId = '',
    this.redirectUri = '',
    this.postLogoutRedirectUri = '',
    this.scopes = const ['openid', 'profile', 'email'],
  });
  factory OidcConfig.generated() {
    final json = jsonDecode(authConfigJson) as Map<String, dynamic>;
    final native = json['native'] as Map<String, dynamic>;
    return OidcConfig(
      mode: json['mode'] as String,
      issuer: json['issuer'] as String,
      clientId: native['client_id'] as String,
      redirectUri: native['redirect_uri'] as String,
      postLogoutRedirectUri: native['post_logout_redirect_uri'] as String,
      scopes: List<String>.from(json['scopes'] as List),
    );
  }
  final String mode, issuer, clientId, redirectUri, postLogoutRedirectUri;
  final List<String> scopes;
  bool get enabled => mode == 'oidc';
  void validate() {
    if (mode == 'none') return;
    final authority = Uri.tryParse(issuer);
    if (mode != 'oidc' ||
        authority == null ||
        authority.scheme != 'https' ||
        authority.host.isEmpty ||
        authority.userInfo.isNotEmpty ||
        authority.hasQuery ||
        authority.hasFragment ||
        clientId.isEmpty ||
        !scopes.contains('openid')) {
      throw const FormatException('Invalid OIDC configuration.');
    }
    for (final value in [redirectUri, postLogoutRedirectUri]) {
      final uri = Uri.tryParse(value);
      if (uri == null ||
          !uri.hasScheme ||
          uri.scheme == 'http' ||
          uri.scheme == 'javascript' ||
          uri.scheme == 'data' ||
          uri.host.isEmpty ||
          uri.path.isEmpty ||
          uri.hasQuery ||
          uri.hasFragment ||
          uri.userInfo.isNotEmpty) {
        throw const FormatException('Invalid OIDC callback.');
      }
    }
    if (redirectUri == postLogoutRedirectUri) {
      throw const FormatException('Callbacks must differ.');
    }
  }

  String get storageKey =>
      'oidc.${base64Url.encode(utf8.encode('$issuer|$clientId'))}';
}
