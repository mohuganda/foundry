import 'package:app_links/app_links.dart';

class DeepLinkService {
  DeepLinkService([AppLinks? appLinks]) : _appLinks = appLinks ?? AppLinks();
  final AppLinks _appLinks;
  Future<Uri?> get initialUri => _appLinks.getInitialLink();
  Stream<Uri> get uriStream => _appLinks.uriLinkStream;
}
