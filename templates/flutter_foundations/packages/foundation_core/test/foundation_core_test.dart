import 'package:foundation_core/foundation_core.dart';
import 'package:test/test.dart';

void main() {
  group('Result', () {
    test('folds a success value', () {
      const result = Success<int>(42);

      final value = result.fold(
        success: (value) => 'value: $value',
        failure: (failure) => failure.message,
      );

      expect(value, 'value: 42');
    });

    test('folds a failure', () {
      const result = Error<int>(ValidationFailure('Invalid input'));

      final value = result.fold(
        success: (value) => 'value: $value',
        failure: (failure) => failure.message,
      );

      expect(value, 'Invalid input');
    });
  });

  group('Paginated', () {
    test('reports when another page is available', () {
      const page = Paginated<int>(items: [1, 2], page: 1, perPage: 2, total: 3);

      expect(page.hasNext, isTrue);
    });

    test('reports when the final page is reached', () {
      const page = Paginated<int>(items: [3], page: 2, perPage: 2, total: 3);

      expect(page.hasNext, isFalse);
    });
  });
}
