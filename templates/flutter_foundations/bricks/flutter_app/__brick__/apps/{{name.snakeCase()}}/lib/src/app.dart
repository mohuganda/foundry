import 'dart:async';

import 'package:flutter/material.dart';
import 'package:foundation_auth/foundation_auth.dart';
import 'package:foundation_identity/foundation_identity.dart';
import 'package:{{name.snakeCase()}}/src/routing/app_router.dart';
import 'package:{{name.snakeCase()}}/src/theme/app_theme.dart';

class AppRoot extends StatefulWidget {
  const AppRoot({this.identity, this.auth, super.key});
  final AppIdentity? identity;
  final OidcAuth? auth;
  @override
  State<AppRoot> createState() => _AppRootState();
}

class _AppRootState extends State<AppRoot> {
  late final OidcAuth authentication;
  @override
  void initState() {
    super.initState();
    authentication = widget.auth ?? OidcAuth(config: OidcConfig.generated());
    unawaited(authentication.initialize());
  }

  @override
  void dispose() {
    if (widget.auth == null) authentication.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final profile = widget.identity ?? appIdentity;
    return MaterialApp.router(
      title: profile.application,
      theme: AppTheme.light(profile),
      darkTheme: AppTheme.dark(profile),
      routerConfig: appRouter,
      builder: (context, child) => ListenableBuilder(
        listenable: authentication,
        builder: (context, _) {
          Widget content;
          switch (authentication.status) {
            case OidcStatus.loading:
              content = IdentitySignIn(identity: profile, pending: true);
            case OidcStatus.error:
              content = IdentitySignIn(
                identity: profile,
                error: authentication.error,
                onSignIn: () => unawaited(authentication.signIn()),
              );
            case OidcStatus.anonymous:
              content = IdentitySignIn(
                identity: profile,
                onSignIn: () => unawaited(authentication.signIn()),
              );
            case OidcStatus.authenticated:
              content = Column(
                children: [
                  TextButton(
                    onPressed: () => unawaited(authentication.signOut()),
                    child: const Text('Sign out'),
                  ),
                  Expanded(child: child ?? const SizedBox.shrink()),
                ],
              );
            case OidcStatus.disabled:
              content = child ?? const SizedBox.shrink();
          }
          return IdentityShell(
            identity: profile,
            authentication:
                authentication.status != OidcStatus.disabled &&
                authentication.status != OidcStatus.authenticated,
            onOpenLink: (link) => unawaited(
              openIdentityLink(context, link, onNavigate: appRouter.go),
            ),
            child: content,
          );
        },
      ),
    );
  }
}
