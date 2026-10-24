import 'package:{{app.snakeCase()}}/features/{{name.snakeCase()}}/domain/models/{{name.snakeCase()}}.dart';
import 'package:{{app.snakeCase()}}/features/{{name.snakeCase()}}/domain/repositories/{{name.snakeCase()}}_repository.dart';

final class {{name.pascalCase()}}RepositoryImpl implements {{name.pascalCase()}}Repository {
  @override
  Future<List<{{name.pascalCase()}}>> list() async {
    // TODO: Connect remote/local data sources.
    return const [];
  }
}
