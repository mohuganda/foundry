import 'package:flutter/widgets.dart';
import 'package:{{name.snakeCase()}}/src/bootstrap.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  runApp(await bootstrap());
}
