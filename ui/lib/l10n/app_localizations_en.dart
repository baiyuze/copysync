// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for English (`en`).
class AppLocalizationsEn extends AppLocalizations {
  AppLocalizationsEn([String locale = 'en']) : super(locale);

  @override
  String get navHistory => 'History';

  @override
  String get navDevices => 'Devices';

  @override
  String get navSettings => 'Settings';

  @override
  String get linkConnecting => 'Connecting…';

  @override
  String get linkServiceNotRunning => 'Service not running';

  @override
  String get linkNoServer => 'No server connection';

  @override
  String linkOnlinePeers(int count) {
    String _temp0 = intl.Intl.pluralLogic(
      count,
      locale: localeName,
      other: '$count devices online',
      one: '1 device online',
    );
    return '$_temp0';
  }

  @override
  String get linkServerConnected => 'Connected to the server';

  @override
  String linkConnectionLost(String detail) {
    return 'Connection lost: $detail';
  }

  @override
  String get linkServiceDisconnected => 'Service disconnected';

  @override
  String get linkServiceDisabled => 'Background sync is off';

  @override
  String get errServiceNotConnected =>
      'Not connected to the background service';

  @override
  String get errUnimplemented => 'This isn\'t available yet';

  @override
  String get errUnavailable => 'Can\'t reach the background service';

  @override
  String errServiceStart(String detail) {
    return 'Couldn\'t start the background service: $detail';
  }

  @override
  String errServiceRestart(String detail) {
    return 'Couldn\'t restart the background service: $detail';
  }

  @override
  String errAutostart(String detail) {
    return 'Couldn\'t set CopySync to start at login: $detail';
  }

  @override
  String get thisMacTitle => 'This Mac';

  @override
  String get thisPcTitle => 'This computer';

  @override
  String get thisMacInText => 'this Mac';

  @override
  String get thisPcInText => 'this computer';

  @override
  String get bothMacs => 'both Macs';

  @override
  String get bothDevices => 'both devices';

  @override
  String get thisDeviceShort => 'This device';

  @override
  String fromDevice(String name) {
    return 'From $name';
  }

  @override
  String get cancel => 'Cancel';

  @override
  String get close => 'Close';

  @override
  String get save => 'Save';

  @override
  String get retry => 'Try again';

  @override
  String get continueAction => 'Continue';

  @override
  String get justNow => 'Just now';

  @override
  String secondsAgo(int n) {
    return '${n}s ago';
  }

  @override
  String minutesAgo(int n) {
    String _temp0 = intl.Intl.pluralLogic(
      n,
      locale: localeName,
      other: '$n minutes ago',
      one: '1 minute ago',
    );
    return '$_temp0';
  }

  @override
  String hoursAgo(int n) {
    String _temp0 = intl.Intl.pluralLogic(
      n,
      locale: localeName,
      other: '$n hours ago',
      one: '1 hour ago',
    );
    return '$_temp0';
  }

  @override
  String daysAgo(int n) {
    String _temp0 = intl.Intl.pluralLogic(
      n,
      locale: localeName,
      other: '$n days ago',
      one: 'Yesterday',
    );
    return '$_temp0';
  }

  @override
  String monthDay(int month, int day) {
    return '$month/$day';
  }

  @override
  String durationHours(int n) {
    String _temp0 = intl.Intl.pluralLogic(
      n,
      locale: localeName,
      other: '$n hours',
      one: '1 hour',
    );
    return '$_temp0';
  }

  @override
  String durationDays(int n) {
    String _temp0 = intl.Intl.pluralLogic(
      n,
      locale: localeName,
      other: '$n days',
      one: '1 day',
    );
    return '$_temp0';
  }

  @override
  String get historyEmptySubtitle =>
      'Copy on any of your devices and it shows up here';

  @override
  String historyCount(int count) {
    String _temp0 = intl.Intl.pluralLogic(
      count,
      locale: localeName,
      other: '$count items',
      one: '1 item',
    );
    return '$_temp0';
  }

  @override
  String historyCountKept(int count, String ttl) {
    String _temp0 = intl.Intl.pluralLogic(
      count,
      locale: localeName,
      other: '$count items',
      one: '1 item',
    );
    return '$_temp0, kept for $ttl';
  }

  @override
  String get filterAll => 'All';

  @override
  String get filterText => 'Text';

  @override
  String get filterImage => 'Images';

  @override
  String get filterFile => 'Files';

  @override
  String get historyClearTooltip => 'Clear history';

  @override
  String get historyEmptyTitle => 'Nothing here yet';

  @override
  String get historyEmptyFilteredTitle => 'Nothing of this kind';

  @override
  String get historyEmptyBody =>
      'Copy text, an image or a file on any paired device and it appears here.';

  @override
  String get clearConfirmTitle => 'Clear the whole history?';

  @override
  String clearConfirmBody(String device) {
    return 'Everything synced from your devices, and the cached files, will be removed from $device. The original files aren\'t affected.';
  }

  @override
  String get clear => 'Clear';

  @override
  String get historyCleared => 'History cleared';

  @override
  String get permDeniedTitle => 'CopySync can\'t read the clipboard';

  @override
  String get permAskTitle => 'Allow CopySync to read the clipboard';

  @override
  String get permDeniedBody =>
      'What you copy on this Mac won\'t be synced. In System Settings, set CopySync Daemon to Allow.';

  @override
  String get permAskBody =>
      'Otherwise macOS asks every time. In System Settings, set CopySync Daemon to Allow.';

  @override
  String get openSystemSettings => 'Open System Settings';

  @override
  String get fetch => 'Pull to this device';

  @override
  String get fetchStarted => 'Pulling…';

  @override
  String get transferFailed => 'Transfer failed';

  @override
  String get expired => 'Expired';

  @override
  String get putOnClipboard => 'Copy to clipboard';

  @override
  String get putOnClipboardDone => 'On the clipboard';

  @override
  String get deleteRecord => 'Delete this item';

  @override
  String get deleted => 'Deleted';

  @override
  String get emptyRecord => '(empty)';

  @override
  String get imageRecord => 'Image';

  @override
  String itemsSummary(String name, int others, int count) {
    String _temp0 = intl.Intl.pluralLogic(
      others,
      locale: localeName,
      other: '$others more',
      one: '1 more',
    );
    return '$name and $_temp0';
  }

  @override
  String get imagePreviewTitle => 'Image preview';

  @override
  String get closePreview => 'Close preview';

  @override
  String get previousImage => 'Previous image';

  @override
  String get nextImage => 'Next image';

  @override
  String get previewHint => 'Scroll or pinch to zoom, drag to move around';

  @override
  String get fitWindow => 'Fit to window';

  @override
  String get previewDeleted => 'This image has been deleted from the history';

  @override
  String get previewExpired => 'This image has expired and can\'t be previewed';

  @override
  String get previewDownloading => 'Downloading the image…';

  @override
  String get previewDownloadFailed => 'The download failed. Try again.';

  @override
  String get previewNotDownloaded => 'This image isn\'t on this computer yet';

  @override
  String get downloadAndPreview => 'Download and preview';

  @override
  String get previewUnavailable => 'This image can\'t be previewed right now';

  @override
  String get imageSemantic => 'Image from the history';

  @override
  String get previewUnreadable =>
      'Can\'t read the image. The file may have been cleaned up or damaged.';

  @override
  String get devicesSubtitleEmpty =>
      'Once paired, your devices share one clipboard';

  @override
  String devicesSubtitle(int paired, int online) {
    return '$paired paired, $online online';
  }

  @override
  String get addDevice => 'Add device';

  @override
  String thisDeviceFootnote(String both) {
    return 'When pairing, $both show the same two lines of security fingerprints, one for each device. Confirm only if they match character for character.';
  }

  @override
  String get noDevicesTitle => 'No paired devices yet';

  @override
  String get noDevicesBody =>
      'Open CopySync on your other computer and link the two with a 6-character pairing code.';

  @override
  String get noDevicesNeedServer =>
      'Connect to a signaling server in Settings before pairing.';

  @override
  String get pairedDevices => 'Paired devices';

  @override
  String get pairedFootnote =>
      'Direct: data goes straight between the two devices. Relay: when the network blocks a direct path, data goes through your server and stays end-to-end encrypted.';

  @override
  String get offline => 'Offline';

  @override
  String get direct => 'Direct';

  @override
  String get relay => 'Relay';

  @override
  String get online => 'Online';

  @override
  String get moreActions => 'More';

  @override
  String get copyFingerprint => 'Copy fingerprint';

  @override
  String get unpairEllipsis => 'Unpair…';

  @override
  String get fingerprintCopied => 'Fingerprint copied';

  @override
  String unpairTitle(String name) {
    return 'Unpair “$name”?';
  }

  @override
  String get unpairBody =>
      'The two devices will stop syncing. To sync again, you\'ll need to pair them and check the fingerprints again.';

  @override
  String get unpair => 'Unpair';

  @override
  String get unpaired => 'Unpaired';

  @override
  String get generateCode => 'Show a code';

  @override
  String get enterCode => 'Enter a code';

  @override
  String get codeFailedTitle => 'Couldn\'t create a pairing code';

  @override
  String get enterOnOther => 'Enter this code on your other device';

  @override
  String get codeCopied => 'Code copied';

  @override
  String get codeHint => 'Click to copy. Valid for 5 minutes.';

  @override
  String afterEnterFootnote(String both) {
    return 'Once it\'s entered, $both show the same two lines of security fingerprints. Confirm only if they match.';
  }

  @override
  String get codeInvalid => 'A pairing code is 6 letters or digits';

  @override
  String get enterCodeShown => 'Enter the code shown on your other device';

  @override
  String get pairingDone => 'Paired';

  @override
  String get pairingRejected => 'Pairing declined';

  @override
  String get verifyTitle => 'Check the security fingerprints';

  @override
  String pairingWith(String name) {
    return 'Pairing with “$name”';
  }

  @override
  String verifyBody(String both) {
    return 'The two lines on $both must be identical. If even one character differs, someone may be tampering with the connection, so decline.';
  }

  @override
  String get rejectMismatch => 'Mismatch, decline';

  @override
  String get confirmMatch => 'They match, pair';

  @override
  String get server => 'Server';

  @override
  String get serverFootnote =>
      'The server only helps devices find each other and set up a direct connection. It can\'t see what they send.';

  @override
  String get signalingUrl => 'Signaling server';

  @override
  String get signalingUrlDesc =>
      'Press Return to save; it reconnects right away';

  @override
  String get signalingUrlPlaceholder => 'ws://server-address:8787/signal';

  @override
  String get publicStun => 'Probe with public servers';

  @override
  String get publicStunDesc =>
      'On networks with several uplinks, such as offices with two ISPs, this finds more direct paths. Public servers only see your public IP.';

  @override
  String get connectionStatus => 'Connection';

  @override
  String get connectedDesc => 'Connected to the signaling server';

  @override
  String get disconnectedDesc =>
      'Can\'t reach the server. Check the address and that the server is running.';

  @override
  String get connected => 'Connected';

  @override
  String get notConnected => 'Not connected';

  @override
  String get deviceName => 'Device name';

  @override
  String get deviceNameDesc => 'Shown in the list on your other devices';

  @override
  String get language => 'Language';

  @override
  String get languageSystem => 'System';

  @override
  String get sync => 'Sync';

  @override
  String get syncFootnote =>
      'Items under the limit are pushed to your other devices as you copy. Larger ones sync as a record; pull them from History when you need them, so big files don\'t waste bandwidth or disk space.';

  @override
  String get autoSyncLimit => 'Auto-sync limit';

  @override
  String get autoSyncLimitDesc => 'Larger files are pulled by hand';

  @override
  String get autoApply => 'Put received items on the clipboard';

  @override
  String get autoApplyDesc =>
      'When off, put them on the clipboard from History yourself';

  @override
  String get whatToSync => 'What to sync';

  @override
  String get syncFiles => 'Files and folders';

  @override
  String get syncText => 'Plain text';

  @override
  String get syncRich => 'Rich text';

  @override
  String get syncImages => 'Images and screenshots';

  @override
  String get storage => 'Storage';

  @override
  String get storageFootnote =>
      'Expired items and cached files are cleaned up automatically.';

  @override
  String get keepHistory => 'Keep history for';

  @override
  String get keepCache => 'Keep cached files for';

  @override
  String get keepCacheDesc => 'Never longer than the history';

  @override
  String get cacheUsage => 'Cache size';

  @override
  String get backgroundSync => 'Background sync';

  @override
  String get backgroundFootnote =>
      'When off, nothing syncs and CopySync doesn\'t start at login. History and pairings are kept; open the app to turn it back on.';

  @override
  String get backgroundService => 'Background service';

  @override
  String get backgroundServiceDesc =>
      'Starts at login and keeps syncing with the window closed';

  @override
  String get restart => 'Restart';

  @override
  String get restarted => 'Background service restarted';

  @override
  String get disableEllipsis => 'Turn off…';

  @override
  String get about => 'About';

  @override
  String get aboutDesc => 'Open source under the AGPL-3.0';

  @override
  String get homepage => 'Website';

  @override
  String get feedback => 'Report a problem';

  @override
  String get disableTitle => 'Turn off background sync?';

  @override
  String disableBody(String device) {
    return 'CopySync on $device will stop syncing and won\'t start at login. History and pairings are kept.';
  }

  @override
  String get disable => 'Turn off';

  @override
  String get disabled => 'Background sync is off';

  @override
  String get urlMustBeWs => 'The address must start with ws:// or wss://';

  @override
  String get serverUpdated => 'Server updated. Reconnecting…';

  @override
  String get deviceNameUpdated => 'Device name updated';

  @override
  String get svcExtractTitle => 'Install or extract CopySync first';

  @override
  String get svcExtractBody =>
      'CopySync is running straight from the zip, and its temporary folder disappears when you close it. Run the installer, or extract the whole folder and open it from there, to turn on background sync.';

  @override
  String get svcMoveTitle => 'Move CopySync to Applications first';

  @override
  String get svcMoveBody =>
      'CopySync is running straight from the disk image. Drag it into Applications, eject the disk image, and open it from Applications to turn on background sync.';

  @override
  String get svcStartTitle => 'Get started with CopySync';

  @override
  String get svcStartBody =>
      'CopySync runs a small background process that notices clipboard changes and stays connected to your other devices. Once on, it starts at login and keeps syncing even with this window closed.';

  @override
  String get svcEnable => 'Turn on background sync';

  @override
  String get svcStaleTitle => 'Turn background sync on again';

  @override
  String get svcStaleBodyWin =>
      'CopySync has moved, but starting at login still points to the old place. Turn it on again so the background process uses this copy. History and pairings are kept.';

  @override
  String get svcStaleBodyMac =>
      'CopySync has moved, or was installed with an old script. Turn it on again so the background process uses this app. History and pairings are kept.';

  @override
  String get svcReenable => 'Turn on again';

  @override
  String get svcStoppedTitle => 'The background service isn\'t running';

  @override
  String get svcStoppedBody =>
      'It may have just quit; the system usually restarts it within seconds. If you stay on this page, restart it with the button below.';

  @override
  String get svcDevBodyWin =>
      'This is a development build with no copysyncd.exe next to it. Run scripts\\build-windows.ps1 in the source folder to build a complete version.';

  @override
  String get svcDevBodyMac =>
      'This is a development build without the background service. Run ./scripts/install-macos.sh in the source folder to install it.';

  @override
  String svcLogPath(String path) {
    return 'Log: $path';
  }

  @override
  String get svcCanDisable =>
      'You can turn background sync off anytime in Settings.';
}
