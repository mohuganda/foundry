import 'dart:io';
import 'package:file_picker/file_picker.dart';
import 'package:image_picker/image_picker.dart';
import 'package:path_provider/path_provider.dart';
import 'package:share_plus/share_plus.dart';
import 'package:url_launcher/url_launcher.dart';

class FileService {
  FileService({ImagePicker? imagePicker})
    : _images = imagePicker ?? ImagePicker();
  final ImagePicker _images;
  Future<PlatformFile?> pickDocument() async =>
      (await FilePicker.platform.pickFiles())?.files.single;
  Future<XFile?> takePhoto() => _images.pickImage(source: ImageSource.camera);
  Future<XFile?> pickImage() => _images.pickImage(source: ImageSource.gallery);
  Future<Directory> appDocumentsDirectory() =>
      getApplicationDocumentsDirectory();
  Future<void> shareFile(String path, {String? text}) async {
    await SharePlus.instance.share(
      ShareParams(files: [XFile(path)], text: text),
    );
  }

  Future<bool> openUrl(Uri uri) =>
      launchUrl(uri, mode: LaunchMode.externalApplication);
}
