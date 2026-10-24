import 'package:{{app.snakeCase()}}/features/{{name.snakeCase()}}/domain/models/{{name.snakeCase()}}.dart';

abstract interface class {{name.pascalCase()}}Repository {
  Future<List<{{name.pascalCase()}}>> list();
}
