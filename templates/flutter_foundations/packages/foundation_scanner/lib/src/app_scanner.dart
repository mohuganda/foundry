import 'package:flutter/material.dart';
import 'package:mobile_scanner/mobile_scanner.dart';

class AppScanner extends StatelessWidget {
  const AppScanner({super.key, required this.onCode});
  final ValueChanged<String> onCode;
  @override
  Widget build(BuildContext context) => MobileScanner(
    onDetect: (capture) {
      for (final code in capture.barcodes) {
        final value = code.rawValue;
        if (value != null && value.isNotEmpty) {
          onCode(value);
          break;
        }
      }
    },
  );
}
