import 'package:drift/drift.dart';

class {{name.pascalCase()}}Rows extends Table {
  TextColumn get id => text()();
  DateTimeColumn get updatedAt => dateTime().withDefault(currentDateAndTime)();
  @override
  Set<Column<Object>> get primaryKey => {id};
}
