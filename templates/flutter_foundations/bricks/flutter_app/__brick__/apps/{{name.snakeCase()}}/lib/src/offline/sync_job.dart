abstract interface class SyncJob {
  String get id;
  Future<void> execute();
}
