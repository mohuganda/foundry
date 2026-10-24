import 'package:sentry_flutter/sentry_flutter.dart';

abstract interface class CrashReporter {
  Future<void> captureException(Object error, StackTrace stackTrace);
}

final class NoopCrashReporter implements CrashReporter {
  const NoopCrashReporter();
  @override
  Future<void> captureException(Object error, StackTrace stackTrace) async {}
}

final class SentryCrashReporter implements CrashReporter {
  const SentryCrashReporter();
  @override
  Future<void> captureException(Object error, StackTrace stackTrace) =>
      Sentry.captureException(error, stackTrace: stackTrace).then((_) {});
}
