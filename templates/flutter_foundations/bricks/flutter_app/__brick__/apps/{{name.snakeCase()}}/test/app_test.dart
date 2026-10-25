import 'package:flutter_test/flutter_test.dart';
import 'package:foundation_auth/foundation_auth.dart';
import 'package:foundation_identity/foundation_identity.dart';
import 'package:{{name.snakeCase()}}/src/app.dart';

class TestAuth extends OidcAuth {
  TestAuth() : super(config: const OidcConfig());
  @override
  Future<void> initialize() async {}
  void transition(OidcStatus next) {
    status = next;
    notifyListeners();
  }
}

void main() {
  testWidgets('starter app uses the configured organisation shell', (
    tester,
  ) async {
    final auth = OidcAuth(config: const OidcConfig());
    addTearDown(auth.dispose);
    await tester.pumpWidget(AppRoot(auth: auth));
    await tester.pumpAndSettle();
    expect(find.byType(IdentityShell), findsOneWidget);
    expect(find.text(appIdentity.organisation), findsWidgets);
  });
  testWidgets('login shell changes with the authentication lifecycle', (
    tester,
  ) async {
    final auth = TestAuth()..transition(OidcStatus.anonymous);
    addTearDown(auth.dispose);
    await tester.pumpWidget(AppRoot(auth: auth));
    await tester.pumpAndSettle();
    expect(
      tester.widget<IdentityShell>(find.byType(IdentityShell)).authentication,
      isTrue,
    );
    expect(find.text('Sign in'), findsOneWidget);
    auth.transition(OidcStatus.authenticated);
    await tester.pumpAndSettle();
    expect(
      tester.widget<IdentityShell>(find.byType(IdentityShell)).authentication,
      isFalse,
    );
    expect(find.text('Sign out'), findsOneWidget);
    auth.transition(OidcStatus.loading);
    await tester.pump();
    expect(
      tester.widget<IdentityShell>(find.byType(IdentityShell)).authentication,
      isTrue,
    );
    expect(find.text('Opening sign-in…'), findsWidgets);
  });
}
