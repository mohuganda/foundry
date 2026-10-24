import 'package:flutter/material.dart';

class {{name.pascalCase()}}Page extends StatelessWidget {
  const {{name.pascalCase()}}Page({super.key});
  @override
  Widget build(BuildContext context) => Scaffold(
    appBar: AppBar(title: const Text('{{name.titleCase()}}')),
    body: const SafeArea(child: Center(child: Text('{{name.titleCase()}}'))),
  );
}
