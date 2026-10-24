import 'package:package_info_plus/package_info_plus.dart';

class AppInfo {
  const AppInfo({
    required this.name,
    required this.version,
    required this.buildNumber,
    required this.packageName,
  });
  final String name, version, buildNumber, packageName;
}

class AppInfoService {
  const AppInfoService();
  Future<AppInfo> load() async {
    final p = await PackageInfo.fromPlatform();
    return AppInfo(
      name: p.appName,
      version: p.version,
      buildNumber: p.buildNumber,
      packageName: p.packageName,
    );
  }
}
