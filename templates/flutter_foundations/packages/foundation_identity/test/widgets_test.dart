import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:foundation_identity/foundation_identity.dart';

AppIdentity profile() => AppIdentity.fromJson({
  'organisation': {
    'name': 'Ministry of Health',
    'country': 'Republic of Uganda',
    'website': 'https://health.go.ug',
  },
  'application': {'name': 'Programme Portal', 'environment': 'production'},
  'navigation': [
    {'label': 'Home', 'href': '/'},
  ],
  'links': {'privacy': '/privacy'},
});

void main() {
  testWidgets('shell displays identity and dispatches public links', (
    tester,
  ) async {
    Uri? opened;
    final identity = profile();
    await tester.pumpWidget(
      MaterialApp(
        theme: identityTheme(identity),
        home: IdentityShell(
          identity: identity,
          onOpenLink: (uri) => opened = uri,
          child: const Center(child: Text('Dashboard')),
        ),
      ),
    );
    expect(find.text('Republic of Uganda'), findsOneWidget);
    expect(find.text('Dashboard'), findsOneWidget);
    expect(find.byType(Image), findsNothing);
    expect(find.text('Organisation website'), findsOneWidget);
    await tester.tap(find.text('Privacy'));
    expect(opened.toString(), '/privacy');
    expect(find.text('production environment'), findsNothing);
  });

  testWidgets('sign-in and error layouts invoke the real callbacks', (
    tester,
  ) async {
    var calls = 0;
    await tester.pumpWidget(
      MaterialApp(
        home: Scaffold(
          body: IdentitySignIn(identity: profile(), onSignIn: () => calls++),
        ),
      ),
    );
    await tester.tap(find.text('Continue'));
    expect(calls, 1);
    await tester.pumpWidget(
      MaterialApp(
        home: Scaffold(
          body: IdentityError(message: 'Offline', onRetry: () => calls++),
        ),
      ),
    );
    await tester.tap(find.text('Try again'));
    expect(calls, 2);
    expect(find.text('Offline'), findsOneWidget);
  });

  testWidgets('shell supports a small display with larger text', (
    tester,
  ) async {
    await tester.binding.setSurfaceSize(const Size(360, 640));
    addTearDown(() => tester.binding.setSurfaceSize(null));
    await tester.pumpWidget(
      MaterialApp(
        builder: (context, child) => MediaQuery(
          data: MediaQuery.of(
            context,
          ).copyWith(textScaler: const TextScaler.linear(1.5)),
          child: child!,
        ),
        home: IdentityShell(identity: profile(), child: const Text('Content')),
      ),
    );
    expect(tester.takeException(), isNull);
  });

  test('application overrides and unsafe links remain isolated', () {
    final original = profile();
    final other = original.withApplication(name: 'Other app');
    expect(original.application, 'Programme Portal');
    expect(other.application, 'Other app');
    expect(identityLink('javascript:alert(1)'), isNull);
    expect(identityLink('//external.example'), isNull);
    expect(identityLink(r'/\external.example'), isNull);
    expect(identityLink('/dashboard').toString(), '/dashboard');
  });
}
