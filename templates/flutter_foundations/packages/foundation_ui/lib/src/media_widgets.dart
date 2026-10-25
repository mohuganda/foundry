import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';
import 'package:lucide_icons_flutter/lucide_icons.dart';
import 'package:shimmer/shimmer.dart';

class AppNetworkImage extends StatelessWidget {
  const AppNetworkImage(
    this.url, {
    super.key,
    this.fit = BoxFit.cover,
    this.width,
    this.height,
  });
  final String url;
  final BoxFit fit;
  final double? width, height;
  @override
  Widget build(BuildContext context) => CachedNetworkImage(
    imageUrl: url,
    fit: fit,
    width: width,
    height: height,
    placeholder: (c, u) => const ShimmerPlaceholder(),
    errorWidget: (c, u, e) => const Icon(LucideIcons.imageOff),
  );
}

class AppSvgAsset extends StatelessWidget {
  const AppSvgAsset(this.asset, {super.key, this.width, this.height});
  final String asset;
  final double? width, height;
  @override
  Widget build(BuildContext context) =>
      SvgPicture.asset(asset, width: width, height: height);
}

class ShimmerPlaceholder extends StatelessWidget {
  const ShimmerPlaceholder({super.key, this.height = 72});
  final double height;
  @override
  Widget build(BuildContext context) => Shimmer.fromColors(
    baseColor: Theme.of(context).colorScheme.surfaceContainerHighest,
    highlightColor: Theme.of(context).colorScheme.surface,
    child: Container(height: height, color: Colors.white),
  );
}
