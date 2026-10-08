#include <flutter/dart_project.h>
#include <flutter/flutter_view_controller.h>
#include <flutter_windows.h>
#include <windows.h>

#include <algorithm>

#include "flutter_window.h"
#include "utils.h"

int APIENTRY wWinMain(_In_ HINSTANCE instance, _In_opt_ HINSTANCE prev,
                      _In_ wchar_t *command_line, _In_ int show_command) {
  // 只开一个窗口：再次打开时把已有的窗口提到前面，与 Mac 上点 Dock 图标一致
  HANDLE single_instance = ::CreateMutexW(nullptr, TRUE, L"Local\\CopySync.UI");
  if (::GetLastError() == ERROR_ALREADY_EXISTS) {
    HWND existing = ::FindWindowW(L"FLUTTER_RUNNER_WIN32_WINDOW", L"CopySync");
    if (existing) {
      if (::IsIconic(existing)) {
        ::ShowWindow(existing, SW_RESTORE);
      }
      ::SetForegroundWindow(existing);
    }
    return EXIT_SUCCESS;
  }

  // Attach to console when present (e.g., 'flutter run') or create a
  // new console when running with a debugger.
  if (!::AttachConsole(ATTACH_PARENT_PROCESS) && ::IsDebuggerPresent()) {
    CreateAndAttachConsole();
  }

  // Initialize COM, so that it is available for use in the library and/or
  // plugins.
  ::CoInitializeEx(nullptr, COINIT_APARTMENTTHREADED);

  flutter::DartProject project(L"data");

  std::vector<std::string> command_line_arguments =
      GetCommandLineArguments();

  project.set_dart_entrypoint_arguments(std::move(command_line_arguments));

  // 默认 960×640，在主屏的工作区里居中，与 Mac 版一致。
  // 之后打开时用上次关闭时的位置与大小（见 FlutterWindow::OnCreate）
  constexpr int kWidth = 960;
  constexpr int kHeight = 640;
  RECT work{};
  ::SystemParametersInfoW(SPI_GETWORKAREA, 0, &work, 0);
  const double scale =
      FlutterDesktopGetDpiForMonitor(::MonitorFromPoint({0, 0}, MONITOR_DEFAULTTOPRIMARY)) / 96.0;
  const double left = (work.left + (work.right - work.left - kWidth * scale) / 2) / scale;
  const double top = (work.top + (work.bottom - work.top - kHeight * scale) / 2) / scale;
  FlutterWindow window(project);
  Win32Window::Point origin(static_cast<unsigned int>(std::max(0.0, left)),
                            static_cast<unsigned int>(std::max(0.0, top)));
  Win32Window::Size size(kWidth, kHeight);
  if (!window.Create(L"CopySync", origin, size)) {
    return EXIT_FAILURE;
  }
  window.SetQuitOnClose(true);

  ::MSG msg;
  while (::GetMessage(&msg, nullptr, 0, 0)) {
    ::TranslateMessage(&msg);
    ::DispatchMessage(&msg);
  }

  ::CoUninitialize();
  ::CloseHandle(single_instance);
  return EXIT_SUCCESS;
}
