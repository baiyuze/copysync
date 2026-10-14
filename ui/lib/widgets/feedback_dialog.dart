import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../app_state.dart';
import '../background_service.dart' show appVersion;
import '../i18n.dart';
import '../icons.dart';
import '../main.dart';
import '../platform.dart';
import '../theme.dart';
import 'common.dart';

/// 新建 issue 的地址：用仓库里的问题模板（.github/ISSUE_TEMPLATE/bug_report.yml），版本号预先填好
const _newIssue = 'https://github.com/baiyuze/pasteferry/issues/new?template=bug_report.yml';

Future<void> showFeedbackDialog(BuildContext context) =>
    showDialog<void>(context: context, builder: (_) => const FeedbackDialog());

/// 反馈问题：先在本机生成诊断报告给用户看，用户自己拷贝或保存，再贴进 GitHub 的 issue。
/// 不替用户上传任何东西：issue 是公开的，贴出去什么由用户看过之后决定。
class FeedbackDialog extends StatefulWidget {
  const FeedbackDialog({super.key});

  @override
  State<FeedbackDialog> createState() => _FeedbackDialogState();
}

class _FeedbackDialogState extends State<FeedbackDialog> {
  String? _report;
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() {
      _report = null;
      _error = null;
    });
    try {
      final report = await AppScope.read(context).exportDiagnostics();
      if (mounted) setState(() => _report = report);
    } catch (e) {
      if (mounted) setState(() => _error = AppState.describeError(e));
    }
  }

  Future<void> _setVerbose(int minutes) async {
    try {
      await AppScope.read(context).setVerboseLog(minutes);
      await _load(); // 报告里有一行写着详细日志开没开
    } catch (e) {
      if (mounted) showToast(context, AppState.describeError(e), error: true);
    }
  }

  Future<void> _copy() async {
    await Clipboard.setData(ClipboardData(text: _report!));
    if (mounted) showToast(context, context.l10n.diagCopied);
  }

  /// 保存的文件可以作为附件拖进 issue，所以多带些日志：粘贴的报告要塞进 issue 正文，有长度限制
  Future<void> _save() async {
    final l = context.l10n;
    try {
      final report = await AppScope.read(context).exportDiagnostics(logLines: 1000);
      final now = DateTime.now();
      String two(int n) => n.toString().padLeft(2, '0');
      final name = 'PasteFerry-diagnostics-${now.year}${two(now.month)}${two(now.day)}'
          '-${two(now.hour)}${two(now.minute)}.txt';
      final file = File('${downloadsDir()}${Platform.pathSeparator}$name');
      await file.parent.create(recursive: true);
      await file.writeAsString(report);
      revealFile(file.path);
      if (mounted) showToast(context, l.diagSaved(name));
    } catch (e) {
      if (mounted) showToast(context, AppState.describeError(e), error: true);
    }
  }

  /// 先把报告拷贝好，打开新建 issue 的页面后直接粘贴。报告生成失败时也照样打开，只是没有报告可贴
  Future<void> _openIssue() async {
    final l = context.l10n;
    if (_report != null) await Clipboard.setData(ClipboardData(text: _report!));
    openUrl('$_newIssue&version=${Uri.encodeQueryComponent(appVersion)}');
    if (mounted && _report != null) showToast(context, l.diagIssueHint);
  }

  @override
  Widget build(BuildContext context) {
    final l = context.l10n;
    final p = context.palette;
    final state = AppScope.of(context);
    final until = state.verboseLogUntil;
    final verbose = until != null && until.isAfter(DateTime.now());
    String two(int n) => n.toString().padLeft(2, '0');

    return Dialog(
      child: SizedBox(
        width: 600,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(Insets.xl, Insets.lg, Insets.md, 0),
              child: Row(
                children: [
                  Expanded(child: Text(l.feedback, style: context.text.titleMedium)),
                  IconAction(icon: AppIcons.xmark, tooltip: l.close, onTap: () => Navigator.pop(context)),
                ],
              ),
            ),
            // 标题与按钮固定，中间的说明、报告、详细日志在窗口矮的时候可以滚动
            Flexible(
              child: SingleChildScrollView(
                padding: const EdgeInsets.fromLTRB(Insets.xl, Insets.md, Insets.xl, 0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Icon(AppIcons.lockShield, size: 15, color: p.textDim),
                        const SizedBox(width: Insets.sm),
                        Expanded(child: Text(l.diagIntro, style: context.text.bodySmall)),
                      ],
                    ),
                    const SizedBox(height: Insets.md),
                    Container(
                      height: 220,
                      decoration: BoxDecoration(
                        color: p.fill,
                        borderRadius: BorderRadius.circular(Radii.sm),
                      ),
                      child: Stack(
                        children: [
                          Positioned.fill(child: _reportView(context)),
                          Positioned(
                            top: 4,
                            right: 4,
                            child: IconAction(icon: AppIcons.arrowClockwise, tooltip: l.diagRefresh, onTap: _load),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: Insets.lg),
                    Row(
                      children: [
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(l.diagVerboseTitle, style: context.text.labelLarge),
                              const SizedBox(height: 2),
                              Text(l.diagVerboseDesc, style: context.text.bodySmall),
                            ],
                          ),
                        ),
                        const SizedBox(width: Insets.lg),
                        if (verbose) ...[
                          Text(l.diagVerboseUntil('${two(until.hour)}:${two(until.minute)}'),
                              style: context.text.bodySmall),
                          TextButton(onPressed: () => _setVerbose(0), child: Text(l.diagVerboseOff)),
                        ] else
                          OutlinedButton(onPressed: () => _setVerbose(30), child: Text(l.diagVerboseOn)),
                      ],
                    ),
                  ],
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(Insets.xl),
              // 与系统对话框一样靠右排，主按钮在最右；放不下时折到下一行
              child: Wrap(
                alignment: WrapAlignment.end,
                spacing: Insets.sm,
                runSpacing: Insets.sm,
                children: [
                  OutlinedButton(onPressed: _report == null ? null : _save, child: Text(l.diagSave)),
                  OutlinedButton(onPressed: _report == null ? null : _copy, child: Text(l.diagCopy)),
                  FilledButton(onPressed: _openIssue, child: Text(l.diagNewIssue)),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _reportView(BuildContext context) {
    final l = context.l10n;
    if (_error != null) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(Insets.lg),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(l.diagFailed(_error!),
                  textAlign: TextAlign.center, style: context.text.bodySmall?.copyWith(color: context.palette.danger)),
              TextButton(onPressed: _load, child: Text(l.retry)),
            ],
          ),
        ),
      );
    }
    if (_report == null) {
      return const Center(child: SizedBox(width: 16, height: 16, child: CircularProgressIndicator(strokeWidth: 2)));
    }
    return Scrollbar(
      child: SingleChildScrollView(
        padding: const EdgeInsets.fromLTRB(Insets.md, Insets.sm, 36, Insets.sm),
        child: SelectableText(
          _report!,
          style: context.text.bodySmall?.copyWith(fontFamily: monoFamily, fontSize: 11, height: 1.45, color: context.palette.text),
        ),
      ),
    );
  }
}
