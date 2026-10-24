import 'package:permission_handler/permission_handler.dart';

class PermissionService {
  const PermissionService();
  Future<bool> requestCamera() => _request(Permission.camera);
  Future<bool> requestLocation() => _request(Permission.locationWhenInUse);
  Future<bool> requestNotifications() => _request(Permission.notification);
  Future<bool> requestPhotos() => _request(Permission.photos);
  Future<bool> _request(Permission permission) async {
    final status = await permission.request();
    return status.isGranted || status.isLimited;
  }

  Future<bool> openSettings() => openAppSettings();
}
