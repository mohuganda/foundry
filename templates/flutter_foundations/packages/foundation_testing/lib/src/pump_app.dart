import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';

extension PumpFoundationApp on WidgetTester {
  Future<void> pumpFoundationApp(
    Widget child, {
    List<dynamic> overrides = const [],
  }) {
    return pumpWidget(
      ProviderScope(
        overrides: [...overrides],
        child: MaterialApp(home: child),
      ),
    );
  }
}
