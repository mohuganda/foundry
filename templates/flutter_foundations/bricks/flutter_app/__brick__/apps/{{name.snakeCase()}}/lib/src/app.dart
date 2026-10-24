import 'package:flutter/material.dart';
import 'package:{{name.snakeCase()}}/src/routing/app_router.dart';
import 'package:{{name.snakeCase()}}/src/theme/app_theme.dart';

class AppRoot extends StatelessWidget {
  const AppRoot({super.key});
  @override
  Widget build(BuildContext context) => MaterialApp.router(
    title: '{{name.titleCase()}}',
    theme: AppTheme.light,
    darkTheme: AppTheme.dark,
    routerConfig: appRouter,
  );
}
