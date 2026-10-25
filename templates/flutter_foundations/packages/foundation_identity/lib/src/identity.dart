import 'dart:convert';

import 'package:foundation_identity/src/generated_identity.dart';

/// Public organisation identity; ownership contacts are not part of this model.
class AppIdentity {
  AppIdentity.fromJson(Map<String, dynamic> data)
    : _data = Map.unmodifiable(
        Map<String, dynamic>.of(data)..remove('ownership'),
      );

  final Map<String, dynamic> _data;

  Map<String, dynamic> section(String key) =>
      Map<String, dynamic>.from(_data[key] as Map? ?? const {});

  String field(String sectionName, String key, [String fallback = '']) =>
      section(sectionName)[key] as String? ?? fallback;

  String get organisation => field('organisation', 'name');
  String get country => field('organisation', 'country');
  String get application => field('application', 'name', 'Application');
  String get description => field('application', 'description');
  String get environment => field('application', 'environment', 'development');
  String get programme => field('application', 'programme');
  String get department => field('application', 'department');
  String get supportEmail => field('application', 'support_email');
  String get website => field('organisation', 'website');

  Map<String, dynamic> get branding => section('branding');
  Map<String, dynamic> get tokens =>
      Map<String, dynamic>.from(branding['tokens'] as Map? ?? const {});
  Map<String, dynamic> get assets =>
      Map<String, dynamic>.from(branding['assets'] as Map? ?? const {});
  List<Map<String, dynamic>> get navigation => [
    for (final item in _data['navigation'] as List? ?? const [])
      Map<String, dynamic>.from(item as Map),
  ];

  String reportHeader(
    String title, {
    String version = '',
    String issuedOn = '',
  }) {
    final document = section('documents');
    final reportVersion = version.isEmpty
        ? document['default_version'] as String? ?? 'draft'
        : version;
    return [
      country,
      organisation,
      application,
      '${document['title_prefix'] ?? ''}$title',
      'Version: $reportVersion',
      if (issuedOn.isNotEmpty) 'Date: $issuedOn',
      if ((document['footer_text'] as String? ?? '').isNotEmpty)
        document['footer_text'] as String,
    ].where((line) => line.isNotEmpty).join('\n');
  }
}

/// Snapshot generated from identity/identity.json; run foundry identity sync to update.
final appIdentity = AppIdentity.fromJson(
  jsonDecode(identityJson) as Map<String, dynamic>,
);
