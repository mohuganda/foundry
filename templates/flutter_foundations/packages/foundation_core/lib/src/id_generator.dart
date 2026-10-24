import 'package:uuid/uuid.dart';

final class IdGenerator {
  const IdGenerator([this._uuid = const Uuid()]);
  final Uuid _uuid;
  String v4() => _uuid.v4();
  String v7() => _uuid.v7();
}
