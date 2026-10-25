import 'dart:async';
import 'package:flutter_test/flutter_test.dart';
import 'package:foundation_auth/foundation_auth.dart';

const config = OidcConfig(
  mode: 'oidc',
  issuer: 'https://auth.example.org/realms/test',
  clientId: 'mobile',
  redirectUri: 'ug.go.health.test://auth/callback',
  postLogoutRedirectUri: 'ug.go.health.test://auth/signed-out',
);
OidcTokens token({bool expired = false}) => OidcTokens(
  accessToken: 'access',
  refreshToken: 'refresh',
  idToken: 'id',
  expiresAt: DateTime.now().add(Duration(hours: expired ? -1 : 1)),
);

class Store implements OidcSessionStore {
  OidcTokens? value;
  int reads = 0;
  @override
  Future<OidcTokens?> read() async {
    reads++;
    return value;
  }

  @override
  Future<void> write(OidcTokens tokens) async {
    value = tokens;
  }

  @override
  Future<void> clear() async {
    value = null;
  }
}

class Platform implements OidcPlatform {
  @override
  bool supported = true;
  int logins = 0, refreshes = 0, logouts = 0;
  bool failRefresh = false;
  Completer<OidcTokens>? pending;
  @override
  Future<OidcTokens> signIn(OidcConfig config) async {
    logins++;
    return token();
  }

  @override
  Future<OidcTokens> refresh(OidcConfig config, OidcTokens tokens) async {
    refreshes++;
    if (failRefresh) throw StateError('sensitive');
    return pending?.future ?? Future.value(token());
  }

  @override
  Future<void> signOut(OidcConfig config, OidcTokens tokens) async {
    logouts++;
  }
}

void main() {
  test(
    'disabled authentication never reads storage or calls the provider',
    () async {
      final store = Store();
      final platform = Platform();
      final auth = OidcAuth(
        config: const OidcConfig(),
        store: store,
        platform: platform,
      );
      addTearDown(auth.dispose);
      await auth.initialize();
      await auth.signIn();
      expect(auth.status, OidcStatus.disabled);
      expect(store.reads, 0);
      expect(platform.logins, 0);
      expect(await auth.accessToken(), isNull);
    },
  );
  test(
    'restores, refreshes once and isolates provider/client storage',
    () async {
      final store = Store()..value = token(expired: true);
      final platform = Platform();
      final auth = OidcAuth(config: config, store: store, platform: platform);
      addTearDown(auth.dispose);
      await auth.initialize();
      expect(auth.status, OidcStatus.authenticated);
      expect(platform.refreshes, 1);
      expect(await auth.accessToken(), 'access');
      expect(
        config.storageKey,
        isNot(
          const OidcConfig(
            issuer: 'https://other.example',
            clientId: 'mobile',
          ).storageKey,
        ),
      );
    },
  );
  test(
    'invalid refresh clears the session and exposes safe feedback',
    () async {
      final store = Store()..value = token(expired: true);
      final platform = Platform()..failRefresh = true;
      final auth = OidcAuth(config: config, store: store, platform: platform);
      addTearDown(auth.dispose);
      await auth.initialize();
      expect(store.value, isNull);
      expect(auth.status, OidcStatus.error);
      expect(auth.error, isNot(contains('sensitive')));
      expect(await auth.accessToken(), isNull);
    },
  );
  test(
    'sign-out during concurrent refresh cannot restore the session',
    () async {
      final store = Store()..value = token();
      final platform = Platform()..pending = Completer<OidcTokens>();
      final auth = OidcAuth(config: config, store: store, platform: platform);
      addTearDown(auth.dispose);
      await auth.initialize();
      final a = auth.refreshAccessToken();
      final b = auth.refreshAccessToken();
      expect(platform.refreshes, 1);
      await auth.signOut();
      platform.pending!.complete(token());
      expect(await a, isNull);
      expect(await b, isNull);
      expect(store.value, isNull);
      expect(auth.status, OidcStatus.anonymous);
      expect(platform.logouts, 1);
    },
  );
  test('unsupported native platforms fail explicitly', () async {
    final auth = OidcAuth(
      config: config,
      store: Store(),
      platform: Platform()..supported = false,
    );
    addTearDown(auth.dispose);
    await auth.initialize();
    expect(auth.status, OidcStatus.error);
    expect(auth.error, contains('Android, iOS and macOS'));
  });
  test('configuration rejects unsafe issuers and callbacks', () {
    expect(config.validate, returnsNormally);
    expect(
      const OidcConfig(
        mode: 'oidc',
        issuer: 'http://auth.example',
        clientId: 'mobile',
      ).validate,
      throwsFormatException,
    );
  });
}
