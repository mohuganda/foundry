import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:{{app.snakeCase()}}/features/{{name.snakeCase()}}/data/repositories/{{name.snakeCase()}}_repository_impl.dart';
import 'package:{{app.snakeCase()}}/features/{{name.snakeCase()}}/domain/models/{{name.snakeCase()}}.dart';
import 'package:{{app.snakeCase()}}/features/{{name.snakeCase()}}/domain/repositories/{{name.snakeCase()}}_repository.dart';

final {{name.camelCase()}}RepositoryProvider = Provider<{{name.pascalCase()}}Repository>(
  (ref) => {{name.pascalCase()}}RepositoryImpl(),
);

final {{name.camelCase()}}ListProvider = FutureProvider<List<{{name.pascalCase()}}>>(
  (ref) => ref.watch({{name.camelCase()}}RepositoryProvider).list(),
);
