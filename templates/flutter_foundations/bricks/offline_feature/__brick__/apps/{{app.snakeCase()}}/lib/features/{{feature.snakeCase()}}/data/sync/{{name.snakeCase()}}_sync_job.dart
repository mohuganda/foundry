import 'package:foundation_sync/foundation_sync.dart';

final class {{name.pascalCase()}}SyncJob implements SyncJob {
  @override
  String get key => '{{feature.snakeCase()}}.{{name.snakeCase()}}';
  @override
  Future<void> run() async {
    /* push pending changes, then pull server changes */
  }
}
