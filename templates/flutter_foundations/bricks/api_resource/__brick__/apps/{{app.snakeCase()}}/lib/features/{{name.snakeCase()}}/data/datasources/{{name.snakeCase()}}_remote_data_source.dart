import 'package:foundation_network/foundation_network.dart';

final class {{name.pascalCase()}}RemoteDataSource {
  const {{name.pascalCase()}}RemoteDataSource(this._client);
  final ApiClient _client;

  Future<List<Map<String, dynamic>>> list() async {
    final response = await _client.dio.get<List<dynamic>>('/{{name.snakeCase()}}');
    return (response.data ?? const []).cast<Map<String, dynamic>>();
  }
}
