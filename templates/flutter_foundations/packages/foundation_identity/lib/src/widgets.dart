import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';
import 'package:foundation_identity/src/identity.dart';
import 'package:foundation_identity/src/login_design.dart';
import 'package:lucide_icons_flutter/lucide_icons.dart';
import 'package:url_launcher/url_launcher.dart';

Color _colour(Map<String, dynamic> palette, String key, Color fallback) {
  final value = palette[key] as String? ?? '';
  final parsed = int.tryParse(value.replaceFirst('#', ''), radix: 16);
  return parsed == null ? fallback : Color(0xff000000 | parsed);
}

/// Theme from configured tokens. Custom font files still need platform registration.
ThemeData identityTheme(
  AppIdentity identity, {
  Brightness brightness = Brightness.light,
}) {
  final tokens = identity.tokens;
  final palette = Map<String, dynamic>.from(
    tokens[brightness == Brightness.dark ? 'dark' : 'light'] as Map? ??
        const {},
  );
  final dark = brightness == Brightness.dark;
  final surface = _colour(
    palette,
    'surface',
    dark ? const Color(0xff192620) : Colors.white,
  );
  final text = _colour(
    palette,
    'text',
    dark ? Colors.white : const Color(0xff182723),
  );
  final focus = _colour(
    palette,
    'focus',
    dark ? Colors.amberAccent : const Color(0xff9c4700),
  );
  final border = _colour(palette, 'border', Colors.grey);
  final radius = (tokens['radius'] as num? ?? 8).toDouble();
  final borderWidth = (tokens['border_width'] as num? ?? 1).toDouble();
  final families =
      (tokens['font_family'] as String? ?? 'Arial, Helvetica, sans-serif')
          .split(',')
          .map((value) => value.replaceAll("'", '').replaceAll('"', '').trim())
          .where((value) => value.isNotEmpty)
          .toList();
  final family = families.isEmpty ? 'Arial' : families.first;
  final base = ThemeData(
    brightness: brightness,
    useMaterial3: true,
    fontFamily: family == 'system-ui' || family == 'sans-serif'
        ? 'Roboto'
        : family,
    fontFamilyFallback: [...families.skip(1), 'Roboto', 'sans-serif'],
    colorScheme: ColorScheme.fromSeed(
      seedColor: _colour(palette, 'primary', Colors.teal),
      brightness: brightness,
      primary: _colour(palette, 'primary', Colors.teal),
      onPrimary: _colour(palette, 'on_primary', Colors.white),
      secondary: _colour(palette, 'success', Colors.green),
      surface: surface,
      onSurface: text,
      onSurfaceVariant: _colour(
        palette,
        'muted',
        dark ? const Color(0xffb7c9bf) : const Color(0xff52645c),
      ),
      error: _colour(
        palette,
        'error',
        dark ? const Color(0xffffb4b4) : const Color(0xffa31515),
      ),
    ),
    scaffoldBackgroundColor: _colour(
      palette,
      'background',
      dark ? const Color(0xff111a17) : const Color(0xfff4f5f1),
    ),
    focusColor: focus,
    dividerColor: border,
    cardTheme: CardThemeData(
      color: surface,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(radius),
        side: BorderSide(color: border, width: borderWidth),
      ),
    ),
    inputDecorationTheme: InputDecorationTheme(
      border: OutlineInputBorder(borderRadius: BorderRadius.circular(radius)),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(radius),
        borderSide: BorderSide(color: focus, width: 2),
      ),
    ),
  );
  return base.copyWith(
    textTheme: base.textTheme.apply(
      fontSizeFactor: (tokens['font_size'] as num? ?? 16) / 16,
    ),
  );
}

class _BrandImage extends StatelessWidget {
  const _BrandImage(this.path, this.label);
  final String path;
  final String label;
  @override
  Widget build(BuildContext context) {
    if (path.isEmpty) return const SizedBox.shrink();
    if (path.endsWith('.svg')) {
      return SvgPicture.asset(
        path,
        package: 'foundation_identity',
        width: 80,
        height: 48,
        semanticsLabel: label,
        errorBuilder: (context, error, stackTrace) => const SizedBox.shrink(),
      );
    }
    return Image.asset(
      path,
      package: 'foundation_identity',
      width: 80,
      height: 48,
      fit: BoxFit.contain,
      semanticLabel: label,
      errorBuilder: (context, error, stackTrace) => const SizedBox.shrink(),
    );
  }
}

class _BrandImages extends StatelessWidget {
  const _BrandImages(this.identity);
  final AppIdentity identity;
  @override
  Widget build(BuildContext context) {
    final assets = identity.assets;
    final dark = Theme.of(context).brightness == Brightness.dark;
    final lightLogo = assets['logo_light'] as String? ?? '';
    final darkLogo = assets['logo_dark'] as String? ?? '';
    return Wrap(
      spacing: 8,
      runSpacing: 8,
      children: [
        _BrandImage(
          dark && darkLogo.isNotEmpty
              ? darkLogo
              : lightLogo.isNotEmpty
              ? lightLogo
              : darkLogo,
          '${identity.organisation} logo',
        ),
        _BrandImage(assets['coat_of_arms'] as String? ?? '', 'Coat of arms'),
      ],
    );
  }
}

/// Shared shell. The application supplies navigation/link handling explicitly.
class IdentityShell extends StatelessWidget {
  const IdentityShell({
    required this.identity,
    required this.child,
    this.onOpenLink,
    this.authentication = false,
    super.key,
  });
  final AppIdentity identity;
  final Widget child;
  final ValueChanged<Uri>? onOpenLink;
  final bool authentication;

  Widget _link(String label, Uri href) => TextButton(
    onPressed: onOpenLink == null ? null : () => onOpenLink!(href),
    child: Text(label),
  );

  Widget _footer(double gap) => Padding(
    padding: EdgeInsets.symmetric(horizontal: gap),
    child: Column(
      children: [
        Text(
          '${identity.organisation} · ${identity.application}',
          textAlign: TextAlign.center,
        ),
        Wrap(
          alignment: WrapAlignment.center,
          children: [
            for (final link in identity.publicLinks.entries)
              _link(link.key, link.value),
          ],
        ),
      ],
    ),
  );

  @override
  Widget build(BuildContext context) {
    final gap = (identity.tokens['space_medium'] as num? ?? 16).toDouble();
    if (authentication) {
      return Scaffold(
        body: SafeArea(
          child: LayoutBuilder(
            builder: (context, constraints) => SingleChildScrollView(
              child: ConstrainedBox(
                constraints: BoxConstraints(minHeight: constraints.maxHeight),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const SizedBox.shrink(),
                    child,
                    Padding(padding: EdgeInsets.all(gap), child: _footer(gap)),
                  ],
                ),
              ),
            ),
          ),
        ),
      );
    }

    return Scaffold(
      body: SafeArea(
        child: Column(
          children: [
            Padding(
              padding: EdgeInsets.all(gap),
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _BrandImages(identity),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        if (identity.country.isNotEmpty) Text(identity.country),
                        Text(
                          identity.organisation,
                          style: Theme.of(context).textTheme.titleMedium,
                        ),
                        Text(identity.application),
                        if (identity.department.isNotEmpty ||
                            identity.programme.isNotEmpty)
                          Text(
                            [
                              identity.department,
                              identity.programme,
                            ].where((value) => value.isNotEmpty).join(' / '),
                          ),
                        if (identity.environment != 'production')
                          Text(
                            '${identity.environment} environment',
                            style: Theme.of(context).textTheme.labelSmall,
                          ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
            if (identity.navigation.isNotEmpty)
              Wrap(
                children: [
                  for (final item in identity.navigation)
                    if (identityLink(item['href'] as String? ?? '')
                        case final Uri href)
                      _link(item['label'] as String? ?? '', href),
                ],
              ),
            const Divider(height: 1),
            Expanded(child: child),
            const Divider(height: 1),
            _footer(gap),
          ],
        ),
      ),
    );
  }
}

final _loginCopy = List<String>.from(jsonDecode(loginCopyJson) as List);

/// Shared client sign-in design; connect the real provider action.
class IdentitySignIn extends StatelessWidget {
  const IdentitySignIn({
    required this.identity,
    this.onSignIn,
    this.pending = false,
    this.error = '',
    super.key,
  });
  final AppIdentity identity;
  final VoidCallback? onSignIn;
  final bool pending;
  final String error;
  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final muted = theme.colorScheme.onSurfaceVariant;
    final radius = (identity.tokens['radius'] as num? ?? 8).toDouble();
    final gap = (identity.tokens['space_medium'] as num? ?? 16).toDouble();
    final border = theme.dividerColor;
    final contextLabel = [
      identity.department,
      identity.programme,
    ].where((v) => v.isNotEmpty).join(' / ');
    return DefaultTextStyle.merge(
      style: theme.textTheme.bodyLarge?.copyWith(height: 1.6),
      child: Center(
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: loginMaxWidth),
            child: Container(
              width: double.infinity,
              clipBehavior: Clip.antiAlias,
              decoration: BoxDecoration(
                color: theme.colorScheme.surface,
                border: Border.all(
                  color: border,
                  width: (identity.tokens['border_width'] as num? ?? 1)
                      .toDouble(),
                ),
                borderRadius: BorderRadius.circular(radius),
              ),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  Container(height: 4, color: theme.colorScheme.primary),
                  Padding(
                    padding: const EdgeInsets.all(loginPadding),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        _BrandImages(identity),
                        if (identity.assets.values.any(
                          (value) => value is String && value.isNotEmpty,
                        ))
                          SizedBox(height: gap),
                        if (identity.country.isNotEmpty)
                          Text(
                            identity.country,
                            style: TextStyle(color: muted, fontSize: 14),
                          ),
                        const SizedBox(height: 8),
                        Text(
                          identity.organisation,
                          style: const TextStyle(fontWeight: FontWeight.w700),
                        ),
                        if (identity.environment != 'production')
                          Padding(
                            padding: const EdgeInsets.only(top: 8),
                            child: Text(
                              '${identity.environment} environment',
                              style: theme.textTheme.labelSmall,
                            ),
                          ),
                      ],
                    ),
                  ),
                  Divider(height: 1, color: border),
                  Padding(
                    padding: const EdgeInsets.all(loginPadding),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      children: [
                        Text(
                          _loginCopy[0].toUpperCase(),
                          style: TextStyle(
                            color: muted,
                            fontSize: 12,
                            fontWeight: FontWeight.w700,
                            letterSpacing: 1,
                          ),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          identity.application,
                          style: const TextStyle(
                            fontSize: 20,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                        if (contextLabel.isNotEmpty)
                          Padding(
                            padding: const EdgeInsets.only(top: 8),
                            child: Text(
                              contextLabel,
                              style: TextStyle(color: muted),
                            ),
                          ),
                        const SizedBox(height: 24),
                        Semantics(
                          header: true,
                          child: Text(
                            _loginCopy[1],
                            style: const TextStyle(
                              fontSize: loginTitleSize,
                              fontWeight: FontWeight.w700,
                              height: 1.2,
                            ),
                          ),
                        ),
                        const SizedBox(height: 8),
                        Text(_loginCopy[2], style: TextStyle(color: muted)),
                        if (error.isNotEmpty)
                          Padding(
                            padding: EdgeInsets.only(top: gap),
                            child: Semantics(
                              liveRegion: true,
                              child: Text(
                                error,
                                style: TextStyle(
                                  color: theme.colorScheme.error,
                                ),
                              ),
                            ),
                          ),
                        SizedBox(height: gap),
                        ConstrainedBox(
                          constraints: const BoxConstraints(
                            minHeight: loginControlHeight,
                          ),
                          child: FilledButton.icon(
                            icon: Icon(
                              pending
                                  ? LucideIcons.loaderCircle
                                  : LucideIcons.logIn,
                              size: 20,
                            ),
                            style: FilledButton.styleFrom(
                              textStyle: const TextStyle(
                                fontSize: 16,
                                fontWeight: FontWeight.w600,
                              ),
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(radius),
                              ),
                            ),
                            onPressed: pending ? null : onSignIn,
                            label: Text(_loginCopy[pending ? 4 : 3]),
                          ),
                        ),
                        if (pending)
                          Semantics(
                            liveRegion: true,
                            child: Padding(
                              padding: const EdgeInsets.only(top: 8),
                              child: Text(_loginCopy[4]),
                            ),
                          ),
                        SizedBox(height: gap),
                        Text(
                          _loginCopy[5],
                          style: TextStyle(
                            color: muted,
                            fontSize: 14,
                            height: 1.5,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

/// Accessible loading presentation.
class IdentityLoading extends StatelessWidget {
  const IdentityLoading({this.message = 'Loading…', super.key});
  final String message;
  @override
  Widget build(BuildContext context) => Center(
    child: Semantics(
      liveRegion: true,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          const CircularProgressIndicator(),
          const SizedBox(height: 16),
          Text(message),
        ],
      ),
    ),
  );
}

/// User-safe error presentation with an optional real retry action.
class IdentityError extends StatelessWidget {
  const IdentityError({
    this.message = 'Please try again.',
    this.onRetry,
    super.key,
  });
  final String message;
  final VoidCallback? onRetry;
  @override
  Widget build(BuildContext context) => Center(
    child: Semantics(
      liveRegion: true,
      child: Card(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(
                'Something went wrong',
                style: Theme.of(context).textTheme.headlineSmall,
              ),
              Text(message),
              const SizedBox(height: 16),
              FilledButton(onPressed: onRetry, child: const Text('Try again')),
            ],
          ),
        ),
      ),
    ),
  );
}

/// Reusable document header; issuedOn and version are explicit report inputs.
class IdentityReportHeader extends StatelessWidget {
  const IdentityReportHeader({
    required this.identity,
    required this.title,
    this.version = '',
    this.issuedOn = '',
    super.key,
  });
  final AppIdentity identity;
  final String title;
  final String version;
  final String issuedOn;
  @override
  Widget build(BuildContext context) => Column(
    crossAxisAlignment: CrossAxisAlignment.start,
    children: [
      _BrandImages(identity),
      SelectableText(
        identity.reportHeader(title, version: version, issuedOn: issuedOn),
      ),
    ],
  );
}

/// Open a validated public link, delegating application paths to its router.
Future<void> openIdentityLink(
  BuildContext context,
  Uri link, {
  required ValueChanged<String> onNavigate,
}) async {
  if (!link.hasScheme) {
    onNavigate(link.toString());
    return;
  }
  try {
    if ((link.scheme == 'https' || link.scheme == 'mailto') &&
        await launchUrl(link, mode: LaunchMode.externalApplication)) {
      return;
    }
  } catch (_) {
    // Present safe feedback below; no URLs or auth details are logged.
  }
  if (context.mounted) {
    ScaffoldMessenger.maybeOf(
      context,
    )?.showSnackBar(const SnackBar(content: Text('Unable to open this link.')));
  }
}
