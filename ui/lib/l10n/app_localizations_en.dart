// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for English (`en`).
class AppLocalizationsEn extends AppLocalizations {
  AppLocalizationsEn([String locale = 'en']) : super(locale);

  @override
  String buildDetails(String uiVersion, String serviceVersion, String build) {
    return 'UI $uiVersion · Service $serviceVersion · Build $build';
  }

  @override
  String get notReceived => 'Not received — receive before pasting';

  @override
  String get receivedPreserved => 'Received; use Copy to clipboard when ready';

  @override
  String get clipboardWriteFailed =>
      'Received, but the clipboard could not be updated';

  @override
  String get receivingStatus => 'Receiving…';

  @override
  String get largeFileNotifications => 'Large file reminders';

  @override
  String get notificationsAllowed =>
      'Notifications are enabled. Progress is shown in History.';

  @override
  String get notificationsNotRequested =>
      'Enable notifications so large files do not wait silently. You can always receive them from History.';

  @override
  String get notificationsDenied =>
      'Notifications are blocked in system settings. Receive pending items from History.';

  @override
  String get notificationsUnknown =>
      'Your desktop controls notification delivery and Do Not Disturb. Receive items from History if no reminder appears.';

  @override
  String get notificationsUnavailable =>
      'System notifications are unavailable. Receive large files from History.';

  @override
  String get enableNotifications => 'Enable reminders';

  @override
  String get notificationSettings => 'Notification settings';

  @override
  String get notifyPendingTitle => 'Receive to paste';

  @override
  String notifyPendingBody(String device, String name, String size) {
    return '$device · $name · $size\nReceive this item before pasting. You can also receive it from History.';
  }

  @override
  String get notifyReceivingTitle => 'Receiving…';

  @override
  String notifyReceivingBody(String name) {
    return '$name\nSee History for progress.';
  }

  @override
  String get notifyReadyTitle => 'Ready to paste';

  @override
  String notifyReadyBody(String name) {
    return '$name\nReceived and placed on the clipboard.';
  }

  @override
  String get notifyReceivedTitle => 'Received';

  @override
  String notifyReceivedBody(String name) {
    return '$name\nYour clipboard was kept. Use Copy to clipboard when ready; this is also available in History.';
  }

  @override
  String get notifyFailedTitle => 'Could not receive';

  @override
  String notifyFailedBody(String name, String error) {
    return '$name\n$error\nTry again from this notification or History.';
  }

  @override
  String get notifyClipboardFailedTitle =>
      'Received; clipboard was not updated';

  @override
  String notifyClipboardFailedBody(String name, String error) {
    return '$name\n$error\nUse Copy to clipboard to try again, here or in History.';
  }

  @override
  String get appName => 'PasteFerry';

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
  String get linkLanMode => 'Looking for devices on this network';

  @override
  String get linkServerUnreachable =>
      'Can\'t reach the server; syncing on this network';

  @override
  String get linkNothingOn =>
      'Local network discovery is off and no server is set';

  @override
  String get linkLanBlocked => 'macOS is blocking the local network';

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
    return 'Couldn\'t set PasteFerry to start at login: $detail';
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
  String get permDeniedTitle => 'PasteFerry can\'t read the clipboard';

  @override
  String get permAskTitle => 'Allow PasteFerry to read the clipboard';

  @override
  String get permDeniedBody =>
      'What you copy on this Mac won\'t be synced. In System Settings, set PasteFerry Daemon to Allow.';

  @override
  String get permAskBody =>
      'Otherwise macOS asks every time. In System Settings, set PasteFerry Daemon to Allow.';

  @override
  String get openSystemSettings => 'Open System Settings';

  @override
  String get fetch => 'Receive';

  @override
  String get fetchStarted => 'Receiving… See History for progress';

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
  String get searchHint => 'Search text and file names';

  @override
  String get searchClear => 'Clear search';

  @override
  String searchResultCount(int count) {
    String _temp0 = intl.Intl.pluralLogic(
      count,
      locale: localeName,
      other: '$count matches',
      one: '1 match',
    );
    return '$_temp0';
  }

  @override
  String searchEmptyTitle(String query) {
    return 'No matches for “$query”';
  }

  @override
  String get searchEmptyBody =>
      'Search looks through copied text and the names of copied files.';

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
      'Open PasteFerry on your other computer and link the two with a 6-character pairing code. On the same network, you can also open Search nearby on both and click the other one.';

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
  String get enterAddressEllipsis => 'Enter its address…';

  @override
  String setAddressTitle(String name) {
    return 'Address of “$name”';
  }

  @override
  String setAddressBody(String name) {
    return 'Use this when the two devices are on different subnets, or both changed addresses. Enter the address shown on “$name” under Add device → Show a code.';
  }

  @override
  String get addressSaved => 'Looking for it at that address';

  @override
  String get lanBlockedTitle => 'macOS is blocking the local network';

  @override
  String lanBlockedBody(String device) {
    return 'Turn on PasteFerry Daemon in System Settings → Privacy & Security → Local Network. Until then, $device can\'t connect to devices on this network.';
  }

  @override
  String get searchNearby => 'Search nearby';

  @override
  String get nearbyTitle => 'Nearby devices';

  @override
  String get nearbySearching => 'Looking for nearby devices…';

  @override
  String get nearbyHint =>
      'Open Search nearby on the other computer too, and it shows up here. Click it to pair.';

  @override
  String nearbyVisible(String device) {
    return 'While this window is open, devices on this network can see the name of $device.';
  }

  @override
  String get nearbyPaired => 'Paired';

  @override
  String nearbyRequesting(String name) {
    return 'Contacting “$name”…';
  }

  @override
  String nearbyNoReply(String name) {
    return '“$name” didn\'t answer. Make sure Search nearby is still open on it.';
  }

  @override
  String get nearbyLanOffTitle => 'Local network discovery is off';

  @override
  String get nearbyLanOffBody =>
      'Turn on “Find devices on this network” in Settings to see nearby devices, or pair with a code instead.';

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
  String get localAddress => 'Address';

  @override
  String get addressHint =>
      'On a different subnet? Enter this address on the other device as well.';

  @override
  String get addressCopied => 'Address copied';

  @override
  String get differentSubnet => 'On a different subnet?';

  @override
  String get peerAddress => 'Address shown under the code';

  @override
  String get pairNotFound =>
      'Couldn\'t find that device. If it\'s on a different subnet, enter the address shown under its code.';

  @override
  String get pairNotFoundAddress =>
      'Nothing answered at that address. Check the address and that the other device is still showing its code. A firewall on that device may also be blocking it.';

  @override
  String get codeWrong => 'That code isn\'t right, or it has expired';

  @override
  String get codeExpired =>
      'This code has expired. Show a new one on the other device.';

  @override
  String get addressInvalid =>
      'That address isn\'t valid. Enter something like 192.168.1.20.';

  @override
  String get pairNoChannel =>
      'Turn on local network discovery or set a server in Settings first.';

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
  String get server => 'Server (optional)';

  @override
  String get serverFootnote =>
      'Only needed to sync between different networks, such as home and office. Devices on the same network connect without it. The server only helps devices find each other and can\'t see what they send.';

  @override
  String get signalingUrl => 'Signaling server';

  @override
  String get signalingUrlDesc =>
      'Press Return to save. Leave it empty to sync on this network only.';

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
  String get notSet => 'Not set';

  @override
  String get notSetDesc => 'Syncing with devices on this network only';

  @override
  String get serverCleared => 'Server removed. Syncing on this network only.';

  @override
  String get localNetwork => 'Local network';

  @override
  String get lanFootnote =>
      'Uses UDP port 8788. If the devices are on different subnets, enter the other device\'s address when pairing.';

  @override
  String get lanDiscovery => 'Find devices on this network';

  @override
  String get lanDiscoveryDesc =>
      'Pair and sync without a server. Others on the network only see an anonymous entry that changes daily, and the device name only while Search nearby is open.';

  @override
  String get macLocalNetwork => 'Local Network access';

  @override
  String get macLocalNetworkDesc =>
      'macOS asks once whether PasteFerry may connect to devices on your network. If you declined, turn on PasteFerry Daemon there.';

  @override
  String get allowLan => 'Windows Firewall';

  @override
  String get allowLanDesc =>
      'Allow this when other devices can\'t reach this one, for example when pairing from another subnet. Asks for administrator approval once.';

  @override
  String get allow => 'Allow';

  @override
  String get allowLanDone => 'Allowed. Other devices can now reach this one.';

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
      'Items under the limit are pushed to your other devices as you copy. Larger ones sync as a record; receive them from History when you need them, so big files don\'t waste bandwidth or disk space.';

  @override
  String get autoSyncLimit => 'Auto-sync limit';

  @override
  String get autoSyncLimitDesc =>
      'Larger files must be received before pasting';

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
      'When off, nothing syncs and PasteFerry doesn\'t start at login. History and pairings are kept; open the app to turn it back on.';

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
    return 'PasteFerry on $device will stop syncing and won\'t start at login. History and pairings are kept.';
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
  String get diagIntro =>
      'This report has the app version, your system, settings, devices and recent logs. Device names, the server address, public IP addresses, your user folder and file names are replaced, and what you copy is never written to the logs. Issues on GitHub are public, so read it through before posting.';

  @override
  String get diagVerboseTitle => 'Detailed logging';

  @override
  String get diagVerboseDesc =>
      'For connection problems: turn it on, make the problem happen again, then come back for the report. It turns itself off after 30 minutes.';

  @override
  String get diagVerboseOn => 'Turn on for 30 minutes';

  @override
  String diagVerboseUntil(String time) {
    return 'On until $time';
  }

  @override
  String get diagVerboseOff => 'Turn off';

  @override
  String get diagRefresh => 'Regenerate';

  @override
  String get diagSave => 'Save as file';

  @override
  String diagSaved(String name) {
    return 'Saved to your Downloads folder: $name';
  }

  @override
  String get diagCopy => 'Copy report';

  @override
  String get diagCopied => 'Report copied';

  @override
  String get diagNewIssue => 'New issue on GitHub';

  @override
  String get diagIssueHint =>
      'Report copied. Paste it into the issue\'s “Diagnostics report” box.';

  @override
  String diagFailed(String detail) {
    return 'Couldn\'t create the report: $detail';
  }

  @override
  String get svcExtractTitle => 'Install or extract PasteFerry first';

  @override
  String get svcExtractBody =>
      'PasteFerry is running straight from the zip, and its temporary folder disappears when you close it. Run the installer, or extract the whole folder and open it from there, to turn on background sync.';

  @override
  String get svcMoveTitle => 'Move PasteFerry to Applications first';

  @override
  String get svcMoveBody =>
      'PasteFerry is running straight from the disk image. Drag it into Applications, eject the disk image, and open it from Applications to turn on background sync.';

  @override
  String get svcStartTitle => 'Get started with PasteFerry';

  @override
  String get svcStartBody =>
      'PasteFerry runs a small background process that notices clipboard changes and stays connected to your other devices. Once on, it starts at login and keeps syncing even with this window closed.';

  @override
  String get svcEnable => 'Turn on background sync';

  @override
  String get svcStaleTitle => 'Turn background sync on again';

  @override
  String get svcStaleBodyWin =>
      'PasteFerry has moved, but starting at login still points to the old place. Turn it on again so the background process uses this copy. History and pairings are kept.';

  @override
  String get svcStaleBodyMac =>
      'PasteFerry has moved, or was installed with an old script. Turn it on again so the background process uses this app. History and pairings are kept.';

  @override
  String get svcReenable => 'Turn on again';

  @override
  String get svcStoppedTitle => 'The background service isn\'t running';

  @override
  String get svcStoppedBody =>
      'It may have just quit; the system usually restarts it within seconds. If you stay on this page, restart it with the button below.';

  @override
  String get svcDevBodyWin =>
      'This is a development build with no pasteferryd.exe next to it. Run scripts\\build-windows.ps1 in the source folder to build a complete version.';

  @override
  String get svcDevBodyLinux =>
      'This is a development build with no pasteferryd next to it. Run ./scripts/package-linux.sh in the source folder to build a complete version.';

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
