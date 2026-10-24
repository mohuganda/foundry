import 'package:flutter/material.dart';

enum AppMessageKind { success, error, warning, info }

extension AppMessengerContext on BuildContext {
  void showAppMessage(
    String message, {
    AppMessageKind kind = AppMessageKind.info,
  }) {
    final scheme = Theme.of(this).colorScheme;
    final background = switch (kind) {
      AppMessageKind.success => scheme.primary,
      AppMessageKind.error => scheme.error,
      AppMessageKind.warning => scheme.tertiary,
      AppMessageKind.info => scheme.inverseSurface,
    };
    ScaffoldMessenger.of(this)
      ..hideCurrentSnackBar()
      ..showSnackBar(
        SnackBar(content: Text(message), backgroundColor: background),
      );
  }

  void success(String message) =>
      showAppMessage(message, kind: AppMessageKind.success);
  void error(String message) =>
      showAppMessage(message, kind: AppMessageKind.error);
  void warning(String message) =>
      showAppMessage(message, kind: AppMessageKind.warning);
  void info(String message) => showAppMessage(message);
}
