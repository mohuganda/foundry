enum SyncStatus { idle, syncing, offline, failed }

class SyncState {
  const SyncState({
    this.status = SyncStatus.idle,
    this.pending = 0,
    this.failed = 0,
    this.lastSyncedAt,
  });
  final SyncStatus status;
  final int pending;
  final int failed;
  final DateTime? lastSyncedAt;
}
