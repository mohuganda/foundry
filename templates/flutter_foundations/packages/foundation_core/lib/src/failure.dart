sealed class Failure {
  const Failure(this.message, {this.cause, this.stackTrace});

  final String message;
  final Object? cause;
  final StackTrace? stackTrace;
}

final class NetworkFailure extends Failure {
  const NetworkFailure(super.message, {super.cause, super.stackTrace});
}

final class UnauthorizedFailure extends Failure {
  const UnauthorizedFailure([super.message = 'Authentication required']);
}

final class ValidationFailure extends Failure {
  const ValidationFailure(super.message, {this.fields = const {}});
  final Map<String, List<String>> fields;
}

final class StorageFailure extends Failure {
  const StorageFailure(super.message, {super.cause, super.stackTrace});
}

final class UnknownFailure extends Failure {
  const UnknownFailure(super.message, {super.cause, super.stackTrace});
}
