import 'package:flutter/material.dart';
import 'package:foundation_identity/foundation_identity.dart';

abstract final class AppTheme {
  static ThemeData light(AppIdentity identity) => identityTheme(identity);
  static ThemeData dark(AppIdentity identity) =>
      identityTheme(identity, brightness: Brightness.dark);
}
