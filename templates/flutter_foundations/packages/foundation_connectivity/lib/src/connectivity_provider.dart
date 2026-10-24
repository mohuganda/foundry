import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:foundation_connectivity/src/connectivity_service.dart';

final connectivityServiceProvider = Provider<ConnectivityService>(
  (ref) => ConnectivityService(),
);
final connectivityProvider = StreamProvider<bool>(
  (ref) => ref.watch(connectivityServiceProvider).changes,
);
