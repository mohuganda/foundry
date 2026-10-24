import 'package:shared_preferences/shared_preferences.dart';

final class KeyValueStore {
  const KeyValueStore(this._preferences);
  final SharedPreferencesAsync _preferences;
  Future<String?> readString(String key) => _preferences.getString(key);
  Future<void> writeString(String key, String value) =>
      _preferences.setString(key, value);
  Future<void> remove(String key) => _preferences.remove(key);
}
