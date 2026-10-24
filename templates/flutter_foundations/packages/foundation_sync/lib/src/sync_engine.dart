import 'dart:async';
import 'package:foundation_connectivity/foundation_connectivity.dart';
import 'package:foundation_observability/foundation_observability.dart';

abstract interface class SyncJob {
  String get key;
  Future<void> run();
}

class SyncEngine {
  SyncEngine({required this.connectivity, required this.logger});
  final ConnectivityService connectivity;
  final AppLogger logger;
  bool _running = false;
  Future<void> run(Iterable<SyncJob> jobs) async {
    if (_running || !await connectivity.isConnected) return;
    _running = true;
    try {
      for (final job in jobs) {
        try {
          logger.info('Sync ${job.key} started');
          await job.run();
          logger.info('Sync ${job.key} completed');
        } catch (e, st) {
          logger.error('Sync ${job.key} failed', e, st);
        }
      }
    } finally {
      _running = false;
    }
  }
}
