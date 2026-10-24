import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:{{app.snakeCase()}}/features/{{name.snakeCase()}}/presentation/providers/{{name.snakeCase()}}_providers.dart';

class {{name.pascalCase()}}Page extends ConsumerWidget {
  const {{name.pascalCase()}}Page({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final items = ref.watch({{name.camelCase()}}ListProvider);
    return Scaffold(
      appBar: AppBar(title: const Text('{{name.titleCase()}}')),
      body: items.when(
        data: (value) => ListView.builder(
          itemCount: value.length,
          itemBuilder: (_, index) => ListTile(title: Text(value[index].id)),
        ),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stackTrace) => Center(child: Text('Error: $error')),
      ),
    );
  }
}
