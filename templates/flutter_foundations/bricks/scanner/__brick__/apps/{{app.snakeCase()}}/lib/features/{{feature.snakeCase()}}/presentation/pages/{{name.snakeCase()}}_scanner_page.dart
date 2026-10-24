import 'package:flutter/material.dart';
import 'package:foundation_scanner/foundation_scanner.dart';

class {{name.pascalCase()}}ScannerPage extends StatelessWidget {
  const {{name.pascalCase()}}ScannerPage({super.key});
  @override
  Widget build(BuildContext context) => Scaffold(
    appBar: AppBar(title: const Text('{{name.titleCase()}} Scanner')),
    body: AppScanner(onCode: (code) => Navigator.of(context).pop(code)),
  );
}
