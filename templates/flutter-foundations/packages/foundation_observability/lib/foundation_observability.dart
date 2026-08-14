abstract interface class AppLogger {
  void info(String message, {Map<String, Object?> context = const {}});
  void error(String message, Object error, {StackTrace? stackTrace});
}
