import 'package:{{name.snakeCase()}}/src/offline/sync_job.dart';

final class SyncEngine {
  SyncEngine(this._jobs);
  final List<SyncJob> _jobs;
  Future<void> syncAll() async {
    for (final job in _jobs) {
      await job.execute();
    }
  }
}
