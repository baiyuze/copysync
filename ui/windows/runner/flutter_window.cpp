#include "flutter_window.h"

#include <flutter_windows.h>

#include <optional>

#include "flutter/generated_plugin_registrant.h"

namespace {

// 窗口的位置与大小存在这里，下次打开时恢复，相当于 Mac 上的 setFrameAutosaveName
constexpr wchar_t kSettingsKey[] = L"Software\\CopySync";
constexpr wchar_t kPlacementValue[] = L"WindowPlacement";

// 窗口最小尺寸（逻辑像素），与 Mac 版一致
constexpr int kMinWidth = 760;
constexpr int kMinHeight = 520;

void SavePlacement(HWND hwnd) {
  WINDOWPLACEMENT placement{sizeof(WINDOWPLACEMENT)};
  if (::GetWindowPlacement(hwnd, &placement)) {
    ::RegSetKeyValueW(HKEY_CURRENT_USER, kSettingsKey, kPlacementValue,
                      REG_BINARY, &placement, sizeof(placement));
  }
}

// 恢复上次的位置与大小。返回上次是否最大化；记录无效（比如那块屏幕已经拔掉了）时不恢复。
bool RestorePlacement(HWND hwnd) {
  WINDOWPLACEMENT placement{};
  DWORD size = sizeof(placement);
  if (::RegGetValueW(HKEY_CURRENT_USER, kSettingsKey, kPlacementValue,
                     RRF_RT_REG_BINARY, nullptr, &placement, &size) != ERROR_SUCCESS ||
      size != sizeof(placement) || placement.length != sizeof(placement)) {
    return false;
  }
  if (!::MonitorFromRect(&placement.rcNormalPosition, MONITOR_DEFAULTTONULL)) {
    return false;
  }
  const bool maximized = placement.showCmd == SW_SHOWMAXIMIZED;
  // 只放好位置，先不显示：窗口等 Flutter 画好第一帧再显示（见 OnCreate）
  placement.showCmd = SW_HIDE;
  ::SetWindowPlacement(hwnd, &placement);
  return maximized;
}

}  // namespace

FlutterWindow::FlutterWindow(const flutter::DartProject& project)
    : project_(project) {}

FlutterWindow::~FlutterWindow() {}

bool FlutterWindow::OnCreate() {
  if (!Win32Window::OnCreate()) {
    return false;
  }
  const bool maximized = RestorePlacement(GetHandle());

  RECT frame = GetClientArea();

  // The size here must match the window dimensions to avoid unnecessary surface
  // creation / destruction in the startup path.
  flutter_controller_ = std::make_unique<flutter::FlutterViewController>(
      frame.right - frame.left, frame.bottom - frame.top, project_);
  // Ensure that basic setup of the controller was successful.
  if (!flutter_controller_->engine() || !flutter_controller_->view()) {
    return false;
  }
  RegisterPlugins(flutter_controller_->engine());
  SetChildContent(flutter_controller_->view()->GetNativeWindow());

  flutter_controller_->engine()->SetNextFrameCallback([this, maximized]() {
    this->Show();
    if (maximized) {
      ::ShowWindow(this->GetHandle(), SW_MAXIMIZE);
    }
  });

  // Flutter can complete the first frame before the "show window" callback is
  // registered. The following call ensures a frame is pending to ensure the
  // window is shown. It is a no-op if the first frame hasn't completed yet.
  flutter_controller_->ForceRedraw();

  return true;
}

void FlutterWindow::OnDestroy() {
  if (flutter_controller_) {
    flutter_controller_ = nullptr;
  }

  Win32Window::OnDestroy();
}

LRESULT
FlutterWindow::MessageHandler(HWND hwnd, UINT const message,
                              WPARAM const wparam,
                              LPARAM const lparam) noexcept {
  // Give Flutter, including plugins, an opportunity to handle window messages.
  if (flutter_controller_) {
    std::optional<LRESULT> result =
        flutter_controller_->HandleTopLevelWindowProc(hwnd, message, wparam,
                                                      lparam);
    if (result) {
      return *result;
    }
  }

  switch (message) {
    case WM_FONTCHANGE:
      flutter_controller_->engine()->ReloadSystemFonts();
      break;
    case WM_GETMINMAXINFO: {
      const double scale =
          FlutterDesktopGetDpiForMonitor(::MonitorFromWindow(hwnd, MONITOR_DEFAULTTONEAREST)) / 96.0;
      auto* info = reinterpret_cast<MINMAXINFO*>(lparam);
      info->ptMinTrackSize.x = static_cast<LONG>(kMinWidth * scale);
      info->ptMinTrackSize.y = static_cast<LONG>(kMinHeight * scale);
      return 0;
    }
    case WM_CLOSE:
      SavePlacement(hwnd);
      break;
  }

  return Win32Window::MessageHandler(hwnd, message, wparam, lparam);
}
