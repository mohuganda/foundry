import 'package:flutter_test/flutter_test.dart';
import 'package:foundation_auth/foundation_auth.dart';
import 'package:foundation_identity/foundation_identity.dart';
import 'package:{{name.snakeCase()}}/src/app.dart';

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
}
