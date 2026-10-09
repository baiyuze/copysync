import 'dart:ui';

import 'package:flutter/widgets.dart';

import 'l10n/app_localizations.dart';

export 'l10n/app_localizations.dart';

/// 界面支持的语言。文案在 tool/l10n/strings.py 里，改完运行它再 flutter gen-l10n。
const supportedLocales = [Locale('en'), Locale('zh'), Locale('ja')];

/// 设置里存的语言（Config.language）→ 界面用的 Locale。空字符串表示跟随系统，返回 null。
Locale? localeFromSetting(String value) => switch (value) {
      'zh-Hans' => const Locale('zh'),
      'en' => const Locale('en'),
      'ja' => const Locale('ja'),
      _ => null,
    };

/// 按系统的语言偏好顺序，取第一个支持的；繁体中文也用简体，都不支持时用英文。
Locale resolveLocale(List<Locale>? preferred) {
  for (final l in preferred ?? const <Locale>[]) {
    for (final s in supportedLocales) {
      if (s.languageCode == l.languageCode) return s;
    }
  }
  return const Locale('en');
}

/// 不在 widget 里的代码（AppState 的错误提示、相对时间等）用它取文案。
/// MaterialApp 每次构建时按当前语言更新；启动前先按系统语言给一个。
AppLocalizations appL10n = lookupAppLocalizations(resolveLocale(PlatformDispatcher.instance.locales));

extension L10nContext on BuildContext {
  AppLocalizations get l10n => AppLocalizations.of(this);
}
