import 'dart:ui';

import 'package:flutter/widgets.dart';

import 'l10n/app_localizations.dart';

export 'l10n/app_localizations.dart';

/// 繁体中文分台湾与香港两种，用词不同（位址与地址、網路与網絡等）。
const zhHantTW = Locale.fromSubtags(languageCode: 'zh', scriptCode: 'Hant', countryCode: 'TW');
const zhHantHK = Locale.fromSubtags(languageCode: 'zh', scriptCode: 'Hant', countryCode: 'HK');

/// 界面支持的语言。文案在 tool/l10n/strings.py 里，改完运行它再 flutter gen-l10n。
const supportedLocales = [Locale('en'), Locale('zh'), zhHantTW, zhHantHK, Locale('ja')];

/// 设置里存的语言（Config.language）→ 界面用的 Locale。空字符串表示跟随系统，返回 null。
Locale? localeFromSetting(String value) => switch (value) {
      'zh-Hans' => const Locale('zh'),
      'zh-Hant-TW' => zhHantTW,
      'zh-Hant-HK' => zhHantHK,
      'en' => const Locale('en'),
      'ja' => const Locale('ja'),
      _ => null,
    };

/// 按系统的语言偏好顺序，取第一个支持的；都不支持时用英文。
Locale resolveLocale(List<Locale>? preferred) {
  for (final l in preferred ?? const <Locale>[]) {
    if (l.languageCode == 'zh') return _chinese(l);
    for (final s in supportedLocales) {
      if (s.languageCode == l.languageCode) return s;
    }
  }
  return const Locale('en');
}

/// 中文用简体还是哪种繁体。各系统报得不一样：macOS 是 zh-Hant-TW，Windows 是 zh-TW，
/// Linux 是 zh_TW，所以写了文字（Hant、Hans）看文字，没写时看地区。
/// 香港、澳门用香港的繁体，其他繁体用台湾的。
Locale _chinese(Locale l) {
  final region = l.countryCode?.toUpperCase();
  final hant = l.scriptCode == null ? const {'TW', 'HK', 'MO'}.contains(region) : l.scriptCode == 'Hant';
  if (!hant) return const Locale('zh');
  return region == 'HK' || region == 'MO' ? zhHantHK : zhHantTW;
}

/// 不在 widget 里的代码（AppState 的错误提示、相对时间等）用它取文案。
/// MaterialApp 每次构建时按当前语言更新；启动前先按系统语言给一个。
AppLocalizations appL10n = lookupAppLocalizations(resolveLocale(PlatformDispatcher.instance.locales));

extension L10nContext on BuildContext {
  AppLocalizations get l10n => AppLocalizations.of(this);
}
