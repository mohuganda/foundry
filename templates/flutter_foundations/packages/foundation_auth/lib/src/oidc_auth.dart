import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:flutter_appauth/flutter_appauth.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:foundation_auth/src/oidc_config.dart';
import 'package:foundation_network/foundation_network.dart';

class OidcTokens {
  factory OidcTokens.fromJson(Map<String, dynamic> json) => OidcTokens(
    accessToken: json['access_token'] as String,
    refreshToken: json['refresh_token'] as String?,
    idToken: json['id_token'] as String?,
    expiresAt: DateTime.parse(json['expires_at'] as String),
  );
  const OidcTokens({
    required this.accessToken,
    required this.expiresAt,
    this.refreshToken,
    this.idToken,
  });
  final String accessToken;
  final String? refreshToken, idToken;
  final DateTime expiresAt;
  bool get needsRefresh =>
      !expiresAt.isAfter(DateTime.now().add(const Duration(seconds: 30)));
  Map<String, dynamic> toJson() => {
    'access_token': accessToken,
    'refresh_token': refreshToken,
    'id_token': idToken,
    'expires_at': expiresAt.toUtc().toIso8601String(),
  };
}

abstract interface class OidcSessionStore {
  Future<OidcTokens?> read();
  Future<void> write(OidcTokens tokens);
  Future<void> clear();
}

class SecureOidcSessionStore implements OidcSessionStore {
  const SecureOidcSessionStore(
    this.key, {
    FlutterSecureStorage storage = const FlutterSecureStorage(),
  }) : _storage = storage;
  final String key;
  final FlutterSecureStorage _storage;
  @override
  Future<OidcTokens?> read() async {
    final value = await _storage.read(key: key);
    return value == null
        ? null
        : OidcTokens.fromJson(jsonDecode(value) as Map<String, dynamic>);
  }

  @override
  Future<void> write(OidcTokens tokens) =>
      _storage.write(key: key, value: jsonEncode(tokens.toJson()));
  @override
  Future<void> clear() => _storage.delete(key: key);
}

abstract interface class OidcPlatform {
  bool get supported;
  Future<OidcTokens> signIn(OidcConfig config);
  Future<OidcTokens> refresh(OidcConfig config, OidcTokens tokens);
  Future<void> signOut(OidcConfig config, OidcTokens tokens);
}

class AppAuthOidcPlatform implements OidcPlatform {
  AppAuthOidcPlatform({FlutterAppAuth? appAuth})
    : _appAuth = appAuth ?? const FlutterAppAuth();
  final FlutterAppAuth _appAuth;
  @override
  bool get supported =>
      !kIsWeb &&
      {
        TargetPlatform.android,
        TargetPlatform.iOS,
        TargetPlatform.macOS,
      }.contains(defaultTargetPlatform);
  OidcTokens _tokens(TokenResponse response, {OidcTokens? previous}) {
    final access = response.accessToken;
    final expiry = response.accessTokenExpirationDateTime;
    if (access == null || access.isEmpty || expiry == null) {
      throw StateError('Incomplete token response.');
    }
    return OidcTokens(
      accessToken: access,
      expiresAt: expiry,
      refreshToken: response.refreshToken ?? previous?.refreshToken,
      idToken: response.idToken ?? previous?.idToken,
    );
  }

  @override
  Future<OidcTokens> signIn(OidcConfig c) async => _tokens(
    await _appAuth.authorizeAndExchangeCode(
      AuthorizationTokenRequest(
        c.clientId,
        c.redirectUri,
        issuer: c.issuer,
        scopes: c.scopes,
      ),
    ),
  );
  @override
  Future<OidcTokens> refresh(OidcConfig c, OidcTokens tokens) async {
    if (tokens.refreshToken == null) throw StateError('Sign-in required.');
    return _tokens(
      await _appAuth.token(
        TokenRequest(
          c.clientId,
          c.redirectUri,
          issuer: c.issuer,
          scopes: c.scopes,
          refreshToken: tokens.refreshToken,
        ),
      ),
      previous: tokens,
    );
  }

  @override
  Future<void> signOut(OidcConfig c, OidcTokens tokens) async {
    await _appAuth.endSession(
      EndSessionRequest(
        issuer: c.issuer,
        idTokenHint: tokens.idToken,
        postLogoutRedirectUrl: c.postLogoutRedirectUri,
      ),
    );
  }
}

enum OidcStatus { disabled, loading, anonymous, authenticated, error }

/// Provider-neutral application boundary; API authorization remains server-side.
class OidcAuth extends ChangeNotifier implements AuthTokenProvider {
  OidcAuth({
    required this.config,
    OidcPlatform? platform,
    OidcSessionStore? store,
  }) : _platform = platform ?? AppAuthOidcPlatform(),
       _store = store ?? SecureOidcSessionStore(config.storageKey),
       status = config.enabled ? OidcStatus.loading : OidcStatus.disabled;
  final OidcConfig config;
  final OidcPlatform _platform;
  final OidcSessionStore _store;
  OidcStatus status;
  String error = '';
  OidcTokens? _tokens;
  Future<void>? _initializing;
  Future<String?>? _refreshing;
  Future<void> _storageTail = Future.value();
  Future<void> _storageAction(Future<void> Function() action) {
    final work = _storageTail.then((_) => action());
    _storageTail = work.catchError((Object _) {});
    return work;
  }

  bool _disposed = false;
  int _generation = 0;
  void _publish(OidcStatus value, [String message = '']) {
    if (_disposed) return;
    status = value;
    error = message;
    notifyListeners();
  }

  Future<void> _fail() async {
    _tokens = null;
    try {
      await _storageAction(_store.clear);
    } catch (_) {
      /* Remain signed out if storage is unavailable. */
    }
    _publish(
      OidcStatus.error,
      'Authentication could not be completed. Please sign in again.',
    );
  }

  Future<void> initialize() => _initializing ??= () async {
    if (!config.enabled) return;
    final generation = _generation;
    try {
      config.validate();
      if (!_platform.supported) {
        _publish(
          OidcStatus.error,
          'OIDC sign-in is supported on Android, iOS and macOS.',
        );
        return;
      }
      final restored = await _store.read();
      if (generation != _generation || _disposed) return;
      _tokens = restored;
      if (_tokens != null) {
        await accessToken();
      } else {
        _publish(OidcStatus.anonymous);
      }
    } catch (_) {
      await _fail();
    }
  }();
  Future<void> signIn() async {
    if (!config.enabled || _disposed) return;
    final generation = ++_generation;
    _publish(OidcStatus.loading);
    try {
      config.validate();
      if (!_platform.supported) throw UnsupportedError('Platform');
      final tokens = await _platform.signIn(config);
      if (generation != _generation || _disposed) return;
      await _storageAction(() => _store.write(tokens));
      if (generation != _generation || _disposed) return;
      _tokens = tokens;
      _publish(OidcStatus.authenticated);
    } on FlutterAppAuthUserCancelledException {
      if (generation == _generation) {
        _tokens = null;
        await _storageAction(_store.clear);
        _publish(OidcStatus.anonymous);
      }
    } catch (_) {
      if (generation == _generation) await _fail();
    }
  }

  Future<void> signOut() async {
    ++_generation;
    final tokens = _tokens;
    _tokens = null;
    try {
      await _storageAction(_store.clear);
      _publish(config.enabled ? OidcStatus.anonymous : OidcStatus.disabled);
      if (tokens != null) await _platform.signOut(config, tokens);
    } catch (_) {
      await _fail();
    }
  }

  @override
  Future<void> clearSession() async {
    ++_generation;
    _tokens = null;
    try {
      await _storageAction(_store.clear);
    } finally {
      _publish(config.enabled ? OidcStatus.anonymous : OidcStatus.disabled);
    }
  }

  @override
  Future<String?> accessToken() async {
    if (!config.enabled || _disposed) return null;
    final tokens = _tokens;
    if (tokens == null) return null;
    if (!tokens.needsRefresh) {
      _publish(OidcStatus.authenticated);
      return tokens.accessToken;
    }
    return refreshAccessToken();
  }

  @override
  Future<String?> refreshAccessToken() {
    if (_refreshing != null) return _refreshing!;
    final tokens = _tokens;
    final generation = _generation;
    if (tokens == null) return Future.value(null);
    return _refreshing = () async {
      try {
        final next = await _platform.refresh(config, tokens);
        if (generation != _generation || _disposed) return null;
        if (next.needsRefresh) throw StateError('Expired session');
        await _storageAction(() => _store.write(next));
        if (generation != _generation || _disposed) return null;
        _tokens = next;
        _publish(OidcStatus.authenticated);
        return next.accessToken;
      } catch (_) {
        if (generation == _generation) await _fail();
        return null;
      } finally {
        _refreshing = null;
      }
    }();
  }

  @override
  void dispose() {
    _disposed = true;
    ++_generation;
    _tokens = null;
    super.dispose();
  }
}
