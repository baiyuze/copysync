import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../app_state.dart';
import '../gen/copysync/v1/daemon.pb.dart';
import '../i18n.dart';
import '../main.dart';
import '../theme.dart';
import 'common.dart';
import '../icons.dart';
import '../platform.dart';

// 保存“添加设备”这条路由。配对请求到达时只移除它，不误关上层的确认框。
final _pairingSetupRoutes = Expando<DialogRoute<void>>();

Future<void> showPairingDialog(BuildContext context) {
  final navigator = Navigator.of(context, rootNavigator: true);
  if (_pairingSetupRoutes[navigator] != null) return Future.value();
  final route = DialogRoute<void>(
    context: context,
    builder: (_) => const PairingDialog(),
  );
  _pairingSetupRoutes[navigator] = route;
  return navigator.push(route).whenComplete(() {
    if (identical(_pairingSetupRoutes[navigator], route)) {
      _pairingSetupRoutes[navigator] = null;
    }
  });
}

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
    final l = context.l10n;
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
                  Expanded(child: Text(l.addDevice, style: context.text.titleMedium)),
                  IconAction(
                    icon: AppIcons.xmark,
                    tooltip: l.close,
                    onTap: () => Navigator.pop(context),
                  ),
                ],
              ),
            ),
            const SizedBox(height: Insets.md),
            Segmented<int>(
              value: _tab,
              segments: [(0, l.generateCode), (1, l.enterCode)],
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
    final l = context.l10n;
    if (_loading) {
      return const Center(child: CircularProgressIndicator(strokeWidth: 2));
    }
    if (_error != null) {
      return Padding(
        padding: const EdgeInsets.all(Insets.xl),
        child: EmptyState(
          icon: AppIcons.exclamationmarkCircle,
          title: l.codeFailedTitle,
          description: _error,
          action: OutlinedButton(onPressed: _generate, child: Text(l.retry)),
        ),
      );
    }

    return _FillOrScroll(
      child: Padding(
        padding: const EdgeInsets.all(Insets.xl),
        child: Column(
          children: [
            Text(l.enterOnOther, style: context.text.bodySmall),
            const SizedBox(height: Insets.lg),
            // 配对码是本页的主角：放大、等宽、字距拉开，方便照着念
            MouseRegion(
              cursor: SystemMouseCursors.click,
              child: GestureDetector(
                onTap: () async {
                  await Clipboard.setData(ClipboardData(text: _code ?? ''));
                  if (context.mounted) showToast(context, l.codeCopied);
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
            Text(l.codeHint, style: context.text.labelSmall),
            const Spacer(),
            Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Icon(AppIcons.lockShield, size: 15, color: p.textDim),
                const SizedBox(width: Insets.sm),
                Expanded(
                  child: Text(
                    l.afterEnterFootnote(Wording.bothDevicesInText),
                    style: context.text.bodySmall,
                  ),
                ),
              ],
            ),
          ],
        ),
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
    if (_busy) return;
    final code = _controller.text.trim();
    if (code.length != 6) {
      setState(() => _error = context.l10n.codeInvalid);
      return;
    }
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      // Shell 统一响应 pendingPairing，生成方和输入方都只有一个确认入口。
      // 事件可能先于 RPC 返回，因此这里也不能 pop 最上层路由。
      await AppScope.read(context).redeemPairingCode(code);
    } catch (e) {
      if (mounted) setState(() => _error = AppState.describeError(e));
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final p = context.palette;
    final l = context.l10n;
    return _FillOrScroll(
      child: Padding(
        padding: const EdgeInsets.all(Insets.xl),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text(l.enterCodeShown, style: context.text.bodySmall, textAlign: TextAlign.center),
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
                  : Text(l.continueAction),
            ),
          ],
        ),
      ),
    );
  }
}

/// 内容放得下时铺满页签的高度，Spacer 照样把脚注推到底部；
/// 放不下时（换成更长的语言、或系统字号调大）改为滚动，而不是溢出。
class _FillOrScroll extends StatelessWidget {
  const _FillOrScroll({required this.child});

  final Widget child;

  @override
  Widget build(BuildContext context) => LayoutBuilder(
        builder: (context, box) => SingleChildScrollView(
          child: ConstrainedBox(
            constraints: BoxConstraints(minHeight: box.maxHeight),
            child: IntrinsicHeight(child: child),
          ),
        ),
      );
}

/// 指纹核对对话框。
///
/// 这是整个安全模型里唯一依赖人的环节：信令服务器能替换转发中的公钥，
/// 但无法让两台设备显示出相同的指纹。所以这一步不能做成"点确定就过"，
/// 必须让用户真的看见并比对。
Future<void> showPairingConfirmDialog(BuildContext context, Device peer) {
  final navigator = Navigator.of(context, rootNavigator: true);
  final setup = _pairingSetupRoutes[navigator];
  if (setup != null) {
    _pairingSetupRoutes[navigator] = null;
    if (setup.isActive) navigator.removeRoute(setup);
  }
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
    if (_busy) return;
    setState(() => _busy = true);
    final state = AppScope.read(context);
    try {
      await state.confirmPairing(widget.peer.pairingSession, accept);
      if (!mounted) return;
      Navigator.pop(context);
      showToast(context, accept ? context.l10n.pairingDone : context.l10n.pairingRejected);
    } catch (e) {
      state.dismissPendingPairing(widget.peer.pairingSession);
      if (!mounted) return;
      Navigator.pop(context);
      showToast(context, AppState.describeError(e), error: true);
    }
  }

  @override
  Widget build(BuildContext context) {
    final p = context.palette;
    final l = context.l10n;
    final peer = widget.peer;
    final self = AppScope.of(context).self;
    final rows = pairingFingerprintRows(self, peer);
    return Dialog(
      child: SizedBox(
        width: 420,
        child: Padding(
          padding: const EdgeInsets.all(Insets.xl),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Text(l.verifyTitle, style: context.text.titleMedium),
              const SizedBox(height: 2),
              Text(l.pairingWith(peer.name), style: context.text.bodySmall),
              const SizedBox(height: Insets.xl),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: Insets.lg, vertical: Insets.md),
                decoration: BoxDecoration(
                  color: p.fill,
                  borderRadius: BorderRadius.circular(Radii.md),
                ),
                child: Column(
                  children: [
                    for (final (name, fingerprint) in rows)
                      Padding(
                        padding: const EdgeInsets.symmetric(vertical: 6),
                        child: Row(
                          children: [
                            Text(
                              fingerprint,
                              style: TextStyle(
                                fontSize: 21,
                                fontWeight: FontWeight.w600,
                                fontFamily: monoFamily,
                                letterSpacing: 1.5,
                                color: p.text,
                              ),
                            ),
                            const SizedBox(width: Insets.md),
                            Expanded(
                              child: Text(
                                name,
                                textAlign: TextAlign.right,
                                maxLines: 1,
                                overflow: TextOverflow.ellipsis,
                                style: context.text.bodySmall,
                              ),
                            ),
                          ],
                        ),
                      ),
                  ],
                ),
              ),
              const SizedBox(height: Insets.lg),
              Text(l.verifyBody(Wording.bothDevicesInText), style: context.text.bodySmall),
              const SizedBox(height: Insets.xl),
              Row(
                children: [
                  Expanded(
                    child: OutlinedButton(
                      onPressed: _busy ? null : () => _respond(false),
                      child: Text(l.rejectMismatch),
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
                          : Text(l.confirmMatch),
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

/// 配对确认时显示的两行：两台设备各自的公钥指纹，按指纹字符串排序。
///
/// 两台 Mac 各自算出的结果完全相同，用户只需比对两块屏幕是否一致。
/// 早先每台只显示「对方」的指纹，两块屏幕上的字符串天然不同，根本无从比对。
///
/// 刻意不把两把公钥合成一串短码：攻击者能同时替换两边看到的公钥，
/// 对一串 60 位的合成码只需生日攻击（约 2^30 次）就能凑出相同结果；
/// 而分别显示两个指纹，攻击者必须为每一把伪造公钥各做一次原像攻击（约 2^60 次）。
List<(String, String)> pairingFingerprintRows(Device? self, Device peer) {
  final rows = <(String, String)>[
    if (self != null) (self.name, self.publicKeyFingerprint),
    (peer.name, peer.publicKeyFingerprint),
  ];
  rows.sort((a, b) => a.$2.compareTo(b.$2));
  return rows;
}
