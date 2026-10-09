import 'dart:io';

import 'package:flutter/material.dart';

import '../app_state.dart';
import '../background_service.dart';
import '../i18n.dart';
import '../main.dart';
import '../theme.dart';
import '../widgets/common.dart';

/// 后台服务没装、没在跑时的引导页。
///
/// 首次打开 App 时用户看到的就是这一页：说明后台进程是做什么的，
/// 然后一个按钮把它注册为登录项。不需要打开终端。
class ServicePage extends StatefulWidget {
  const ServicePage({super.key});

  @override
  State<ServicePage> createState() => _ServicePageState();
}

class _ServicePageState extends State<ServicePage> {
  bool _busy = false;
  String? _error;

  Future<void> _run(Future<void> Function() action) async {
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await action();
    } catch (e) {
      if (mounted) setState(() => _error = AppState.describeError(e));
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final state = AppScope.of(context);
    final windows = Platform.isWindows;
    final l = context.l10n;
    final (title, body, button) = switch (state.serviceState) {
      ServiceState.mustMove when windows => (l.svcExtractTitle, l.svcExtractBody, null),
      ServiceState.mustMove => (l.svcMoveTitle, l.svcMoveBody, null),
      ServiceState.notInstalled => (l.svcStartTitle, l.svcStartBody, l.svcEnable),
      ServiceState.stale => (
          l.svcStaleTitle,
          windows ? l.svcStaleBodyWin : l.svcStaleBodyMac,
          l.svcReenable,
        ),
      ServiceState.stopped => (l.svcStoppedTitle, l.svcStoppedBody, l.restart),
      _ => (l.svcStoppedTitle, windows ? l.svcDevBodyWin : l.svcDevBodyMac, null),
    };

    return Center(
      child: ConstrainedBox(
        constraints: const BoxConstraints(maxWidth: 420),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const AppMark(size: 64),
            const SizedBox(height: Insets.xl),
            Text(title, style: context.text.headlineSmall),
            const SizedBox(height: Insets.sm),
            Text(body, style: context.text.bodyMedium?.copyWith(color: context.palette.textDim)),
            if (button != null) ...[
              const SizedBox(height: Insets.xl),
              FilledButton(
                onPressed: _busy ? null : () => _run(state.enableService),
                child: _busy
                    ? const SizedBox(
                        width: 14,
                        height: 14,
                        child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
                      )
                    : Text(button),
              ),
            ],
            if (_error != null) ...[
              const SizedBox(height: Insets.md),
              Text(_error!, style: context.text.bodySmall?.copyWith(color: context.palette.danger)),
            ],
            const SizedBox(height: Insets.xl),
            Text(
              state.serviceState == ServiceState.stopped
                  ? l.svcLogPath(state.service.displayLogPath)
                  : l.svcCanDisable,
              style: context.text.bodySmall,
            ),
          ],
        ),
      ),
    );
  }
}
