import 'package:flutter/material.dart';

import '../app_state.dart';
import '../background_service.dart';
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
    final (title, body, button) = switch (state.serviceState) {
      ServiceState.mustMove => (
          '先把 CopySync 移到「应用程序」',
          'CopySync 现在是直接从安装盘里打开的。把它拖进「应用程序」文件夹，'
              '推出安装盘，再从「应用程序」里打开，就可以启用后台同步。',
          null,
        ),
      ServiceState.notInstalled => (
          '开始使用 CopySync',
          'CopySync 会在后台运行一个小进程，用来感知剪贴板的变化，并与你的其他设备保持连接。'
              '启用后它随系统登录自动启动，关掉这个窗口也会照常同步。',
          '启用后台同步',
        ),
      ServiceState.stale => (
          '重新启用后台同步',
          'CopySync 换了位置，或者之前用旧版安装脚本装过。重新启用一次，'
              '让后台进程指向当前这个 App。历史记录与配对关系都会保留。',
          '重新启用',
        ),
      ServiceState.stopped => (
          '后台服务没有运行',
          '它可能刚刚退出，系统通常会在几秒内把它重新拉起。'
              '如果一直停在这一页，点下面的按钮重新启动。',
          '重新启动',
        ),
      _ => (
          '后台服务没有运行',
          '这是开发构建，App 里没有内置后台服务。'
              '在源码目录运行 ./scripts/install-macos.sh 安装它。',
          null,
        ),
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
                  ? '运行日志：${state.service.logPath.replaceFirst(RegExp(r'^/Users/[^/]+'), '~')}'
                  : '随时可以在「设置」中停用后台同步。',
              style: context.text.bodySmall,
            ),
          ],
        ),
      ),
    );
  }
}
