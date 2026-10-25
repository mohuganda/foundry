import 'package:flutter_test/flutter_test.dart';
import 'package:foundation_identity/foundation_identity.dart';

void main() {
  test(
    'creates report identity with explicit date/version and optional fields',
    () {
      final identity = AppIdentity.fromJson({
        'organisation': {
          'name': 'Ministry of Health',
          'country': 'Republic of Uganda',
        },
        'application': {'name': 'Programme Portal'},
        'documents': {'title_prefix': 'Report: ', 'default_version': 'draft'},
        'ownership': {
          'operational_contacts': ['private-reference'],
        },
      });
      expect(identity.application, 'Programme Portal');
      expect(identity.section('ownership'), isEmpty);
      expect(
        identity.reportHeader('Summary', version: '2', issuedOn: '2026-10-25'),
        contains('Report: Summary\nVersion: 2\nDate: 2026-10-25'),
      );
    },
  );
}
