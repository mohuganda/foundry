import 'dart:async';

import 'package:dio/dio.dart';
import 'package:foundation_network/src/auth_token_provider.dart';

final class ApiClient {
  ApiClient({
    required String baseUrl,
    AuthTokenProvider? tokenProvider,
    bool enableLogs = false,
  }) : _tokenProvider = tokenProvider,
       dio = Dio(
         BaseOptions(
           baseUrl: baseUrl,
           connectTimeout: const Duration(seconds: 20),
           receiveTimeout: const Duration(seconds: 30),
         ),
       ) {
    dio.interceptors.add(
      InterceptorsWrapper(onRequest: _onRequest, onError: _onError),
    );
    if (enableLogs) {
      dio.interceptors.add(
        LogInterceptor(requestBody: true, responseBody: false),
      );
    }
  }

  final Dio dio;
  final AuthTokenProvider? _tokenProvider;
  Future<String?>? _refreshing;

  Future<void> _onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    final token = await _tokenProvider?.accessToken();
    if (token != null && token.isNotEmpty) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    handler.next(options);
  }

  Future<void> _onError(
    DioException error,
    ErrorInterceptorHandler handler,
  ) async {
    if (error.response?.statusCode != 401 ||
        _tokenProvider == null ||
        error.requestOptions.extra['retried'] == true) {
      handler.next(error);
      return;
    }
    try {
      final token = await (_refreshing ??= _tokenProvider.refreshAccessToken())
          .whenComplete(() => _refreshing = null);
      if (token == null || token.isEmpty) {
        await _tokenProvider.clearSession();
        handler.next(error);
        return;
      }
      final request = error.requestOptions;
      request.extra['retried'] = true;
      request.headers['Authorization'] = 'Bearer $token';
      final response = await dio.fetch<dynamic>(request);
      handler.resolve(response);
    } catch (_) {
      await _tokenProvider.clearSession();
      handler.next(error);
    }
  }
}
