import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../i18n.dart';
import '../main.dart';
import '../theme.dart';
import '../widgets/common.dart';

class OfficialServicePage extends StatelessWidget {
  const OfficialServicePage({super.key});
  String _mb(int bytes) => '${(bytes / 1000000).toStringAsFixed(2)} MB';
  String _time(int seconds, {bool beijing = false}) {
    var d = DateTime.fromMillisecondsSinceEpoch(seconds * 1000, isUtc: true);
    d = beijing ? d.add(const Duration(hours: 8)) : d.toLocal();
    String pad(int n) => n.toString().padLeft(2, '0');
    return '${d.year}-${pad(d.month)}-${pad(d.day)} ${pad(d.hour)}:${pad(d.minute)}';
  }

  @override
  Widget build(BuildContext context) {
    final state = AppScope.of(context), l = context.l10n, p = context.palette;
    final u = state.officialUsage;
    final known = u != null && u.updatedAtUnix.toInt() > 0;
    final fresh =
        known &&
        u.available &&
        DateTime.now().millisecondsSinceEpoch ~/ 1000 -
                u.updatedAtUnix.toInt() <
            60;
    final used = known ? u.uploadBytes.toInt() + u.downloadBytes.toInt() : 0;
    final limit = known ? u.limitBytes.toInt() : 0;
    Widget row(String title, String value) => Padding(
      padding: const EdgeInsets.symmetric(vertical: 10),
      child: Row(
        children: [
          Expanded(
            child: Text(
              title,
              style: context.text.bodyMedium?.copyWith(color: p.textDim),
            ),
          ),
          SelectableText(value, style: context.text.bodyMedium),
        ],
      ),
    );
    final reason = u?.reason ?? '';
    final status = switch (reason) {
      'quota_exhausted' => l.officialQuotaTitle,
      'peer_quota_exhausted' => l.officialPeerQuota,
      'global_exhausted' ||
      'global_paused' ||
      'relay_blocked' ||
      'device_banned' => l.officialPaused,
      _ =>
        !known
            ? l.officialUnused
            : !fresh
            ? l.officialUnavailable
            : state.status?.signalingConnected == true
            ? l.connected
            : l.notConnected,
    };
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        PageHeader(title: l.officialService),
        Expanded(
          child: ListView(
            padding: const EdgeInsets.fromLTRB(
              Insets.xl,
              0,
              Insets.xl,
              Insets.xxl,
            ),
            children: [
              GroupSection(
                title: l.connectionStatus,
                child: GroupBox(
                  children: [
                    Padding(
                      padding: const EdgeInsets.all(Insets.lg),
                      child: Text(
                        status,
                        style: context.text.bodyMedium?.copyWith(
                          color: fresh ? p.text : p.warning,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
              GroupSection(
                title: l.officialToday,
                footnote: l.officialAccounting,
                child: GroupBox(
                  children: [
                    Padding(
                      padding: const EdgeInsets.all(Insets.lg),
                      child: Column(
                        children: [
                          if (known) ...[
                            row(l.officialUsed, _mb(used)),
                            LinearProgressIndicator(
                              value: limit > 0
                                  ? (used / limit).clamp(0.0, 1.0)
                                  : 0,
                              color: used >= limit * .9 ? p.warning : p.accent,
                              backgroundColor: p.fill,
                            ),
                            const SizedBox(height: 12),
                            row(
                              l.officialRemaining,
                              _mb(u.remainingBytes.toInt()),
                            ),
                            row(l.officialDailyLimit, _mb(limit)),
                            row(l.officialUpload, _mb(u.uploadBytes.toInt())),
                            row(
                              l.officialDownload,
                              _mb(u.downloadBytes.toInt()),
                            ),
                            row(
                              l.officialReset,
                              _time(u.resetAtUnix.toInt(), beijing: true),
                            ),
                            row(
                              l.officialUpdated,
                              _time(u.updatedAtUnix.toInt()),
                            ),
                            if (limit > 0 && used >= limit * .9)
                              Padding(
                                padding: const EdgeInsets.only(top: 12),
                                child: Text(
                                  l.officialWarning,
                                  style: context.text.bodySmall?.copyWith(
                                    color: p.warning,
                                  ),
                                ),
                              ),
                          ] else
                            Padding(
                              padding: const EdgeInsets.symmetric(vertical: 24),
                              child: Text(
                                l.officialUnused,
                                style: context.text.bodyMedium,
                              ),
                            ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
              GroupSection(
                title: l.officialDeviceID,
                child: GroupBox(
                  children: [
                    Padding(
                      padding: const EdgeInsets.all(Insets.lg),
                      child: Row(
                        children: [
                          Expanded(
                            child: SelectableText(
                              state.status?.deviceId ?? state.self?.id ?? '',
                              style: context.text.bodyMedium?.copyWith(
                                fontFamily: monoFamily,
                              ),
                            ),
                          ),
                          IconButton(
                            tooltip: MaterialLocalizations.of(
                              context,
                            ).copyButtonLabel,
                            onPressed: () => Clipboard.setData(
                              ClipboardData(
                                text:
                                    state.status?.deviceId ??
                                    state.self?.id ??
                                    '',
                              ),
                            ),
                            icon: const Icon(Icons.copy_outlined, size: 18),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
              Text(
                l.officialPrivacy,
                style: context.text.bodySmall?.copyWith(color: p.textDim),
              ),
            ],
          ),
        ),
      ],
    );
  }
}
