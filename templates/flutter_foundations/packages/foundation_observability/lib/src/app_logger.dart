import 'package:talker_flutter/talker_flutter.dart';

abstract interface class AppLogger {
  void debug(Object message);
  void info(Object message);
  void warning(Object message);
  void error(Object message, [Object? error, StackTrace? stackTrace]);
}

final class TalkerAppLogger implements AppLogger {
  TalkerAppLogger([Talker? talker]) : talker = talker ?? Talker();
  final Talker talker;
  @override
  void debug(Object message) => talker.debug(message);
  @override
  void info(Object message) => talker.info(message);
  @override
  void warning(Object message) => talker.warning(message);
  @override
  void error(Object message, [Object? error, StackTrace? stackTrace]) =>
      talker.error(message, error, stackTrace);
}
