import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';

abstract final class FoundationLocalization {
  static const delegates = [
    GlobalMaterialLocalizations.delegate,
    GlobalWidgetsLocalizations.delegate,
    GlobalCupertinoLocalizations.delegate,
  ];
  static const supportedLocales = [Locale('en')];
}
