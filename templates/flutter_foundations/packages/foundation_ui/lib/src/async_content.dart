import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:foundation_ui/src/app_states.dart';

final class AsyncContent<T> extends StatelessWidget {
  const AsyncContent({
    super.key,
    required this.value,
    required this.builder,
    this.isEmpty,
    this.onRetry,
  });
  final AsyncValue<T> value;
  final Widget Function(T data) builder;
  final bool Function(T data)? isEmpty;
  final VoidCallback? onRetry;

  @override
  Widget build(BuildContext context) => value.when(
    data: (data) =>
        (isEmpty?.call(data) ?? false) ? const AppEmptyState() : builder(data),
    loading: () => const AppLoadingState(),
    error: (error, stack) =>
        AppErrorState(message: error.toString(), onRetry: onRetry),
  );
}
