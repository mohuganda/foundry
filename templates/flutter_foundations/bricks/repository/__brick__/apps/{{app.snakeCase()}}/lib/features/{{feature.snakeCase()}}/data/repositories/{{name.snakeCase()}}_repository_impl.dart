import 'package:{{app.snakeCase()}}/features/{{feature.snakeCase()}}/domain/repositories/{{name.snakeCase()}}_repository.dart';

final class {{name.pascalCase()}}RepositoryImpl implements {{name.pascalCase()}}Repository {
  @override
  Future<void> refresh() async {}
}
