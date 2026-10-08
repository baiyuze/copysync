import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../app_state.dart';
import '../gen/copysync/v1/daemon.pb.dart';
import '../main.dart';
import '../theme.dart';
import 'common.dart';

/// 配对对话框：既能生成配对码，也能输入对方的配对码。
///
/// 两种角色放在一个对话框里切换，因为用户在打开它之前
/// 往往还没想好"我这边生成还是那边生成"。
class PairingDialog extends StatefulWidget {
  const PairingDialog({super.key});

  @override
  State<PairingDialog> createState() => _PairingDialogState();
}

class _PairingDialogState extends State<PairingDialog> {
  int _tab = 0;

  @override
  Widget build(BuildContext context) {
    return Dialog(
      child: SizedBox(
        width: 420,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(Insets.xl, Insets.lg, Insets.md, 0),
              child: Row(
                children: [
                  Expanded(child: Text('添加设备', style: context.text.titleMedium)),
                  IconAction(
                    icon: CupertinoIcons.xmark,
                    tooltip: '关闭',
                    onTap: () => Navigator.pop(context),
                  ),
                ],
              ),
            ),
            const SizedBox(height: Insets.md),
            Segmented<int>(
              value: _tab,
              segments: const [(0, '生成配对码'), (1, '输入配对码')],
              onChanged: (t) => setState(() => _tab = t),
            ),
            SizedBox(
              height: 272,
              // IndexedStack 让两页都常驻：来回切换不会重新生成配对码
              child: IndexedStack(
                index: _tab,
                children: const [_GenerateTab(), _JoinTab()],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _GenerateTab extends StatefulWidget {
  const _GenerateTab();

  @override
  State<_GenerateTab> createState() => _GenerateTabState();
}

class _GenerateTabState extends State<_GenerateTab> {
  String? _code;
  String? _error;
  bool _loading = false;

  @override
  void initState() {
    super.initState();
    _generate();
  }

  Future<void> _generate() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final code = await AppScope.read(context).createPairingCode();
      if (mounted) setState(() => _code = code);
    } catch (e) {
      if (mounted) setState(() => _error = AppState.describeError(e));
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final p = context.palette;
    if (_loading) {
      return const Center(child: CircularProgressIndicator(strokeWidth: 2));
    }
    if (_error != null) {
      return Padding(
        padding: const EdgeInsets.all(Insets.xl),
        child: EmptyState(
          icon: CupertinoIcons.exclamationmark_circle,
          title: '没能生成配对码',
          description: _error,
          action: OutlinedButton(onPressed: _generate, child: const Text('重试')),
        ),
      );
    }

    return Padding(
      padding: const EdgeInsets.all(Insets.xl),
      child: Column(
        children: [
          Text('在另一台设备上输入这个配对码', style: context.text.bodySmall),
          const SizedBox(height: Insets.lg),
          // 配对码是本页的主角：放大、等宽、字距拉开，方便照着念
          MouseRegion(
            cursor: SystemMouseCursors.click,
            child: GestureDetector(
              onTap: () async {
                await Clipboard.setData(ClipboardData(text: _code ?? ''));
                if (context.mounted) showToast(context, '配对码已拷贝');
              },
              child: Container(
                width: double.infinity,
                padding: const EdgeInsets.symmetric(vertical: Insets.lg),
                decoration: BoxDecoration(
                  color: p.fill,
                  borderRadius: BorderRadius.circular(Radii.md),
                ),
                child: Text(
                  _code ?? '',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 34,
                    fontWeight: FontWeight.w600,
                    fontFamily: monoFamily,
                    letterSpacing: 10,
                    color: p.text,
                  ),
                ),
              ),
            ),
          ),
          const SizedBox(height: Insets.sm),
          Text('点击可拷贝，5 分钟内有效', style: context.text.labelSmall),
          const Spacer(),
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Icon(CupertinoIcons.lock_shield, size: 15, color: p.textDim),
              const SizedBox(width: Insets.sm),
              Expanded(
                child: Text(
                  '对方输入后，两台设备会各自显示一串安全指纹，逐字核对一致再确认。',
                  style: context.text.bodySmall,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _JoinTab extends StatefulWidget {
  const _JoinTab();

  @override
  State<_JoinTab> createState() => _JoinTabState();
}

class _JoinTabState extends State<_JoinTab> {
  final _controller = TextEditingController();
  bool _busy = false;
  String? _error;

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    final code = _controller.text.trim();
    if (code.length != 6) {
      setState(() => _error = '配对码是 6 位字母或数字');
      return;
    }
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      final peer = await AppScope.read(context).redeemPairingCode(code);
      if (!mounted) return;
      Navigator.pop(context);
      await showPairingConfirmDialog(context, peer);
    } catch (e) {
      if (mounted) setState(() => _error = AppState.describeError(e));
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final p = context.palette;
    return Padding(
      padding: const EdgeInsets.all(Insets.xl),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text('输入另一台设备上显示的配对码',
              style: context.text.bodySmall, textAlign: TextAlign.center),
          const SizedBox(height: Insets.lg),
          TextField(
            controller: _controller,
            maxLength: 6,
            textCapitalization: TextCapitalization.characters,
            textAlign: TextAlign.center,
            style: const TextStyle(
              fontSize: 28,
              fontWeight: FontWeight.w600,
              fontFamily: monoFamily,
              letterSpacing: 10,
            ),
            inputFormatters: [
              // 配对码字符集不含易混字符，这里统一转大写并过滤
              FilteringTextInputFormatter.allow(RegExp('[a-zA-Z0-9]')),
              TextInputFormatter.withFunction(
                (_, next) => next.copyWith(text: next.text.toUpperCase()),
              ),
            ],
            decoration: InputDecoration(
              counterText: '',
              hintText: '······',
              hintStyle: TextStyle(letterSpacing: 10, fontSize: 28, color: p.textFaint),
              contentPadding: const EdgeInsets.symmetric(vertical: Insets.md),
            ),
            onSubmitted: (_) => _submit(),
          ),
          if (_error != null) ...[
            const SizedBox(height: Insets.sm),
            Text(
              _error!,
              style: context.text.bodySmall?.copyWith(color: p.danger),
              textAlign: TextAlign.center,
            ),
          ],
          const Spacer(),
          FilledButton(
            onPressed: _busy ? null : _submit,
            child: _busy
                ? const SizedBox(
                    width: 14,
                    height: 14,
                    child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
                  )
                : const Text('继续'),
          ),
        ],
      ),
    );
  }
}

/// 指纹核对对话框。
///
/// 这是整个安全模型里唯一依赖人的环节：信令服务器能替换转发中的公钥，
/// 但无法让两台设备显示出相同的指纹。所以这一步不能做成"点确定就过"，
/// 必须让用户真的看见并比对。
Future<void> showPairingConfirmDialog(BuildContext context, Device peer) {
  return showDialog<void>(
    context: context,
    barrierDismissible: false,
    builder: (ctx) => PairingConfirmDialog(peer: peer),
  );
}

class PairingConfirmDialog extends StatefulWidget {
  const PairingConfirmDialog({super.key, required this.peer});

  final Device peer;

  @override
  State<PairingConfirmDialog> createState() => _PairingConfirmDialogState();
}

class _PairingConfirmDialogState extends State<PairingConfirmDialog> {
  bool _busy = false;

  Future<void> _respond(bool accept) async {
    setState(() => _busy = true);
    final state = AppScope.read(context);
    try {
      await state.confirmPairing(widget.peer.pairingSession, accept);
      if (!mounted) return;
      Navigator.pop(context);
      showToast(context, accept ? '配对完成' : '已拒绝配对');
    } catch (e) {
      if (!mounted) return;
      Navigator.pop(context);
      showToast(context, AppState.describeError(e), error: true);
    }
  }

  @override
  Widget build(BuildContext context) {
    final p = context.palette;
    final peer = widget.peer;
    return Dialog(
      child: SizedBox(
        width: 420,
        child: Padding(
          padding: const EdgeInsets.all(Insets.xl),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Text('核对安全指纹', style: context.text.titleMedium),
              const SizedBox(height: 2),
              Text('正在与「${peer.name}」配对', style: context.text.bodySmall),
              const SizedBox(height: Insets.xl),
              Container(
                padding: const EdgeInsets.symmetric(vertical: Insets.lg),
                decoration: BoxDecoration(
                  color: p.fill,
                  borderRadius: BorderRadius.circular(Radii.md),
                ),
                child: Text(
                  peer.publicKeyFingerprint,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 24,
                    fontWeight: FontWeight.w600,
                    fontFamily: monoFamily,
                    letterSpacing: 2,
                    color: p.text,
                  ),
                ),
              ),
              const SizedBox(height: Insets.lg),
              Text(
                '确认这串字符与对方屏幕上显示的完全一致。'
                '不一致说明连接可能被第三方篡改，请拒绝。',
                style: context.text.bodySmall,
              ),
              const SizedBox(height: Insets.xl),
              Row(
                children: [
                  Expanded(
                    child: OutlinedButton(
                      onPressed: _busy ? null : () => _respond(false),
                      child: const Text('不一致，拒绝'),
                    ),
                  ),
                  const SizedBox(width: Insets.md),
                  Expanded(
                    child: FilledButton(
                      onPressed: _busy ? null : () => _respond(true),
                      child: _busy
                          ? const SizedBox(
                              width: 14,
                              height: 14,
                              child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
                            )
                          : const Text('一致，确认配对'),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}
