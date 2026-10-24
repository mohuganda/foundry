import 'package:dio/dio.dart';
import 'package:foundation_core/foundation_core.dart';

Failure mapDioException(DioException error) {
  final status = error.response?.statusCode;
  if (status == 401 || status == 403) {
    return const UnauthorizedFailure();
  }
  if (status == 422) {
    return ValidationFailure(error.message ?? 'Validation failed');
  }
  return NetworkFailure(
    error.message ?? 'Network request failed',
    cause: error,
    stackTrace: error.stackTrace,
  );
}
