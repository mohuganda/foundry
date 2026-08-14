import 'package:flutter/material.dart';

class FoundationPage extends StatelessWidget {
  const FoundationPage({required this.child, super.key});
  final Widget child;

  @override
  Widget build(BuildContext context) => SafeArea(child: child);
}
