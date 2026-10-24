import 'package:flutter/material.dart';
import 'package:flutter_html/flutter_html.dart';
import 'package:flutter_markdown_plus/flutter_markdown_plus.dart';
import 'package:pdfrx/pdfrx.dart';

class AppMarkdownView extends StatelessWidget {
  const AppMarkdownView(this.data, {super.key});
  final String data;
  @override
  Widget build(BuildContext context) => Markdown(data: data, selectable: true);
}

class AppHtmlView extends StatelessWidget {
  const AppHtmlView(this.data, {super.key});
  final String data;
  @override
  Widget build(BuildContext context) => Html(data: data);
}

class AppPdfAssetView extends StatelessWidget {
  const AppPdfAssetView(this.asset, {super.key});
  final String asset;
  @override
  Widget build(BuildContext context) => PdfViewer.asset(asset);
}

class AppPdfUriView extends StatelessWidget {
  const AppPdfUriView(this.uri, {super.key});
  final Uri uri;
  @override
  Widget build(BuildContext context) => PdfViewer.uri(uri);
}
