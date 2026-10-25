import 'dart:async';

import 'package:flutter/material.dart';
import 'package:foundation_identity/foundation_identity.dart';
import 'package:{{name.snakeCase()}}/src/routing/app_router.dart';
import 'package:{{name.snakeCase()}}/src/theme/app_theme.dart';

class AppRoot extends StatelessWidget {
  const AppRoot({this.identity, super.key});
  final AppIdentity? identity;

  @override
  Widget build(BuildContext context) {
    final profile = identity ?? appIdentity;
    return MaterialApp.router(
      title: profile.application,
      theme: AppTheme.light(profile),
      darkTheme: AppTheme.dark(profile),
      routerConfig: appRouter,
      builder: (context, child) => IdentityShell(
        identity: profile,
        onOpenLink: (link) => unawaited(
          openIdentityLink(context, link, onNavigate: appRouter.go),
        ),
        child: child ?? const SizedBox.shrink(),
      ),
    );
  }
}
