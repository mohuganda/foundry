import 'package:flutter/material.dart';

abstract final class AppBreakpoints {
  static const mobile = 600.0;
  static const desktop = 1024.0;
}

class ResponsiveLayout extends StatelessWidget {
  const ResponsiveLayout({
    super.key,
    required this.mobile,
    this.tablet,
    this.desktop,
  });
  final Widget mobile;
  final Widget? tablet;
  final Widget? desktop;
  @override
  Widget build(BuildContext context) => LayoutBuilder(
    builder: (context, c) {
      if (c.maxWidth >= AppBreakpoints.desktop) {
        return desktop ?? tablet ?? mobile;
      }
      if (c.maxWidth >= AppBreakpoints.mobile) {
        return tablet ?? mobile;
      }
      return mobile;
    },
  );
}
