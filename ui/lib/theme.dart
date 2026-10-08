import 'package:flutter/material.dart';

/// CopySync 的设计系统。
///
/// 这是一个常驻后台的系统工具，界面要像 macOS 自带的应用一样安静：
/// 大面积中性的石墨灰，强调色只出现在按钮、选中态和链接上。
/// 不用渐变、不用彩色图标底、不用发光——这些都会让工具显得像一个营销页面。
abstract final class AppColors {
  // 强调色：比系统蓝略深，白底上的文字对比度达到 WCAG AA
  static const accent = Color(0xFF0A66D8);
  static const accentDark = Color(0xFF3B8BEB);
}

/// 间距阶梯。全应用只用这几个值，避免出现 7px、13px 这类随手写的数字，
/// 那是界面显得"不够整齐"的主要来源。
abstract final class Insets {
  static const xs = 4.0;
  static const sm = 8.0;
  static const md = 12.0;
  static const lg = 16.0;
  static const xl = 24.0;
  static const xxl = 32.0;
}

/// 圆角按层级区分：控件小、分组容器中、对话框大。
abstract final class Radii {
  static const sm = 6.0;
  static const md = 10.0;
  static const lg = 14.0;
}

/// 语义化的颜色角色。页面代码只用这些名字，不判断亮暗模式。
@immutable
class AppPalette extends ThemeExtension<AppPalette> {
  const AppPalette({
    required this.accent,
    required this.background,
    required this.sidebar,
    required this.surface,
    required this.border,
    required this.separator,
    required this.text,
    required this.textDim,
    required this.textFaint,
    required this.selection,
    required this.hover,
    required this.fill,
    required this.online,
    required this.warning,
    required this.danger,
  });

  final Color accent;
  final Color background; // 内容区底色
  final Color sidebar;
  final Color surface; // 分组容器、对话框
  final Color border; // 容器外框
  final Color separator; // 容器内的行分隔线
  final Color text;
  final Color textDim;
  final Color textFaint;
  final Color selection; // 侧边栏选中项
  final Color hover;
  final Color fill; // 分段控件轨道、代码块等浅色填充
  final Color online;
  final Color warning;
  final Color danger;

  static const light = AppPalette(
    accent: AppColors.accent,
    background: Color(0xFFF5F5F4),
    sidebar: Color(0xFFEBEBE9),
    surface: Color(0xFFFFFFFF),
    border: Color(0xFFE0E0DD),
    separator: Color(0xFFECECEA),
    text: Color(0xFF1D1D1F),
    textDim: Color(0xFF6E6E73),
    textFaint: Color(0xFFA1A1A6),
    selection: Color(0xFFDCDCD9),
    hover: Color(0x0A000000),
    fill: Color(0xFFEEEEEC),
    // 状态色取 Apple 的深色变体：浅底上作为文字也清晰可读
    online: Color(0xFF248A3D),
    warning: Color(0xFFB25000),
    danger: Color(0xFFD70015),
  );

  static const dark = AppPalette(
    accent: AppColors.accentDark,
    background: Color(0xFF1C1C1E),
    sidebar: Color(0xFF232325),
    surface: Color(0xFF2A2A2C),
    border: Color(0xFF3A3A3C),
    separator: Color(0xFF333335),
    text: Color(0xFFF5F5F7),
    textDim: Color(0xFF98989D),
    textFaint: Color(0xFF636366),
    selection: Color(0xFF3A3A3C),
    hover: Color(0x0FFFFFFF),
    fill: Color(0xFF323234),
    online: Color(0xFF32D74B),
    warning: Color(0xFFFF9F0A),
    danger: Color(0xFFFF453A),
  );

  @override
  AppPalette copyWith() => this;

  @override
  AppPalette lerp(AppPalette? other, double t) {
    if (other == null) return this;
    Color l(Color a, Color b) => Color.lerp(a, b, t)!;
    return AppPalette(
      accent: l(accent, other.accent),
      background: l(background, other.background),
      sidebar: l(sidebar, other.sidebar),
      surface: l(surface, other.surface),
      border: l(border, other.border),
      separator: l(separator, other.separator),
      text: l(text, other.text),
      textDim: l(textDim, other.textDim),
      textFaint: l(textFaint, other.textFaint),
      selection: l(selection, other.selection),
      hover: l(hover, other.hover),
      fill: l(fill, other.fill),
      online: l(online, other.online),
      warning: l(warning, other.warning),
      danger: l(danger, other.danger),
    );
  }
}

extension PaletteAccess on BuildContext {
  AppPalette get palette => Theme.of(this).extension<AppPalette>()!;
  TextTheme get text => Theme.of(this).textTheme;
}

/// 等宽字体：配对码、指纹、服务器地址这类需要逐字符核对的内容。
const monoFamily = 'Menlo';

/// [fontFamily] 只给截图渲染用：App 本身不指定，跟随系统字体（SF Pro 与苹方）。
ThemeData buildTheme(Brightness brightness, {String? fontFamily, List<String>? fontFamilyFallback}) {
  final dark = brightness == Brightness.dark;
  final p = dark ? AppPalette.dark : AppPalette.light;

  final scheme = ColorScheme.fromSeed(seedColor: p.accent, brightness: brightness).copyWith(
    primary: p.accent,
    onPrimary: Colors.white,
    surface: p.background,
    onSurface: p.text,
    error: p.danger,
  );

  // 字号偏小、字重克制：桌面工具的信息密度应当高于移动应用。
  // 层级靠字号与字重区分，不靠颜色。
  var textTheme = TextTheme(
    headlineSmall: TextStyle(
        fontSize: 22, fontWeight: FontWeight.w600, color: p.text, letterSpacing: -0.4),
    titleMedium: TextStyle(
        fontSize: 15, fontWeight: FontWeight.w600, color: p.text, letterSpacing: -0.2),
    titleSmall: TextStyle(fontSize: 13, fontWeight: FontWeight.w600, color: p.text),
    // 输入框默认用 bodyLarge，不定义的话会是移动端的 16px，比周围的字大一号
    bodyLarge: TextStyle(fontSize: 13, color: p.text),
    bodyMedium: TextStyle(fontSize: 13, color: p.text, height: 1.45),
    bodySmall: TextStyle(fontSize: 12, color: p.textDim, height: 1.4),
    labelLarge: TextStyle(fontSize: 13, fontWeight: FontWeight.w500, color: p.text),
    labelMedium: TextStyle(fontSize: 12, fontWeight: FontWeight.w500, color: p.textDim),
    labelSmall: TextStyle(fontSize: 11, fontWeight: FontWeight.w500, color: p.textDim),
  );
  if (fontFamily != null) {
    textTheme = textTheme.apply(fontFamily: fontFamily, fontFamilyFallback: fontFamilyFallback);
  }

  final controlShape = RoundedRectangleBorder(borderRadius: BorderRadius.circular(Radii.sm));
  const controlPadding = EdgeInsets.symmetric(horizontal: 14, vertical: 9);

  return ThemeData(
    useMaterial3: true,
    brightness: brightness,
    colorScheme: scheme,
    scaffoldBackgroundColor: p.background,
    canvasColor: p.surface,
    textTheme: textTheme,
    fontFamily: fontFamily,
    fontFamilyFallback: fontFamilyFallback,
    extensions: [p],
    splashFactory: NoSplash.splashFactory, // 桌面端没有触摸涟漪
    highlightColor: Colors.transparent,
    hoverColor: p.hover,
    dividerTheme: DividerThemeData(color: p.separator, thickness: 1, space: 1),
    iconTheme: IconThemeData(color: p.textDim, size: 16),
    filledButtonTheme: FilledButtonThemeData(
      style: FilledButton.styleFrom(
        backgroundColor: p.accent,
        foregroundColor: Colors.white,
        disabledBackgroundColor: p.fill,
        disabledForegroundColor: p.textFaint,
        shape: controlShape,
        padding: controlPadding,
        minimumSize: const Size(0, 32),
        textStyle: textTheme.labelLarge,
        elevation: 0,
      ),
    ),
    outlinedButtonTheme: OutlinedButtonThemeData(
      style: OutlinedButton.styleFrom(
        foregroundColor: p.text,
        backgroundColor: p.surface,
        shape: controlShape,
        padding: controlPadding,
        minimumSize: const Size(0, 32),
        textStyle: textTheme.labelLarge,
        side: BorderSide(color: p.border),
      ),
    ),
    textButtonTheme: TextButtonThemeData(
      style: TextButton.styleFrom(
        foregroundColor: p.accent,
        shape: controlShape,
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
        minimumSize: const Size(0, 30),
        textStyle: textTheme.labelLarge,
      ),
    ),
    inputDecorationTheme: InputDecorationTheme(
      isDense: true,
      filled: true,
      fillColor: p.surface,
      contentPadding: const EdgeInsets.symmetric(horizontal: 10, vertical: 9),
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(Radii.sm),
        borderSide: BorderSide(color: p.border),
      ),
      enabledBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(Radii.sm),
        borderSide: BorderSide(color: p.border),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(Radii.sm),
        borderSide: BorderSide(color: p.accent, width: 1.5),
      ),
      hintStyle: textTheme.bodyMedium?.copyWith(color: p.textFaint),
    ),
    switchTheme: SwitchThemeData(
      thumbColor: const WidgetStatePropertyAll(Colors.white),
      trackColor: WidgetStateProperty.resolveWith(
          (s) => s.contains(WidgetState.selected) ? p.accent : p.selection),
      trackOutlineColor: const WidgetStatePropertyAll(Colors.transparent),
      materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
    ),
    // Material 3 的对话框与菜单默认没有阴影，只靠颜色区分层级；
    // macOS 的弹窗和菜单都有明显的投影，这里补上
    dialogTheme: DialogThemeData(
      backgroundColor: p.surface,
      surfaceTintColor: Colors.transparent,
      elevation: 18,
      shadowColor: Color(dark ? 0x99000000 : 0x40000000),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(Radii.lg)),
      titleTextStyle: textTheme.titleMedium,
      contentTextStyle: textTheme.bodyMedium?.copyWith(color: p.textDim),
    ),
    popupMenuTheme: PopupMenuThemeData(
      color: p.surface,
      surfaceTintColor: Colors.transparent,
      elevation: 8,
      shadowColor: Color(dark ? 0x99000000 : 0x33000000),
      textStyle: textTheme.bodyMedium,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(Radii.md),
        side: BorderSide(color: p.border),
      ),
    ),
    snackBarTheme: SnackBarThemeData(
      behavior: SnackBarBehavior.floating,
      backgroundColor: dark ? const Color(0xFF3A3A3C) : const Color(0xFF2C2C2E),
      contentTextStyle: const TextStyle(fontSize: 13, color: Colors.white),
      elevation: 6,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(Radii.md)),
    ),
    tooltipTheme: TooltipThemeData(
      decoration: BoxDecoration(
        color: dark ? const Color(0xFF3A3A3C) : const Color(0xFF2C2C2E),
        borderRadius: BorderRadius.circular(Radii.sm),
      ),
      textStyle: const TextStyle(fontSize: 12, color: Colors.white),
      waitDuration: const Duration(milliseconds: 500),
    ),
    progressIndicatorTheme: ProgressIndicatorThemeData(
      color: p.accent,
      linearTrackColor: p.fill,
      circularTrackColor: Colors.transparent,
    ),
    scrollbarTheme: ScrollbarThemeData(
      thickness: const WidgetStatePropertyAll(6),
      radius: const Radius.circular(3),
      thumbColor: WidgetStatePropertyAll(p.textFaint.withValues(alpha: 0.5)),
    ),
  );
}
