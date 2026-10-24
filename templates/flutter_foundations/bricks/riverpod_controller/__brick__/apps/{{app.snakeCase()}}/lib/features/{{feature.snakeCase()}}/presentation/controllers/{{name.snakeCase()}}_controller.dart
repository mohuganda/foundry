import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';

final {{name.camelCase()}}ControllerProvider =
    AsyncNotifierProvider<{{name.pascalCase()}}Controller, void>({{name.pascalCase()}}Controller.new);

class {{name.pascalCase()}}Controller extends AsyncNotifier<void> {
  @override
  FutureOr<void> build() {}

  Future<void> refresh() async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(() async {});
  }
}
