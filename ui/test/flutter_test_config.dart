import 'dart:async';

import 'package:copysync_ui/i18n.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_test/flutter_test.dart';

/// 测试默认按简体中文跑：断言里写的都是中文文案。其他语言在 i18n_test.dart 里单独测。
Future<void> testExecutable(FutureOr<void> Function() testMain) async {
  TestWidgetsFlutterBinding.ensureInitialized().platformDispatcher.localesTestValue = const [Locale('zh')];
  appL10n = lookupAppLocalizations(const Locale('zh'));
  await testMain();
}
