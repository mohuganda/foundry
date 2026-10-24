import 'dart:io';
import 'package:device_info_plus/device_info_plus.dart';

class DeviceInfoService {
  DeviceInfoService([DeviceInfoPlugin? plugin])
    : _plugin = plugin ?? DeviceInfoPlugin();
  final DeviceInfoPlugin _plugin;
  Future<Map<String, Object?>> load() async {
    if (Platform.isAndroid) {
      final i = await _plugin.androidInfo;
      return {
        'platform': 'android',
        'model': i.model,
        'manufacturer': i.manufacturer,
        'sdk': i.version.sdkInt,
      };
    }
    if (Platform.isIOS) {
      final i = await _plugin.iosInfo;
      return {
        'platform': 'ios',
        'model': i.utsname.machine,
        'systemVersion': i.systemVersion,
      };
    }
    return {'platform': Platform.operatingSystem};
  }
}
