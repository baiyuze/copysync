import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:intl/intl.dart' as intl;

import 'app_localizations_en.dart';
import 'app_localizations_ja.dart';
import 'app_localizations_zh.dart';

// ignore_for_file: type=lint

/// Callers can lookup localized strings with an instance of AppLocalizations
/// returned by `AppLocalizations.of(context)`.
///
/// Applications need to include `AppLocalizations.delegate()` in their app's
/// `localizationDelegates` list, and the locales they support in the app's
/// `supportedLocales` list. For example:
///
/// ```dart
/// import 'l10n/app_localizations.dart';
///
/// return MaterialApp(
///   localizationsDelegates: AppLocalizations.localizationsDelegates,
///   supportedLocales: AppLocalizations.supportedLocales,
///   home: MyApplicationHome(),
/// );
/// ```
///
/// ## Update pubspec.yaml
///
/// Please make sure to update your pubspec.yaml to include the following
/// packages:
///
/// ```yaml
/// dependencies:
///   # Internationalization support.
///   flutter_localizations:
///     sdk: flutter
///   intl: any # Use the pinned version from flutter_localizations
///
///   # Rest of dependencies
/// ```
///
/// ## iOS Applications
///
/// iOS applications define key application metadata, including supported
/// locales, in an Info.plist file that is built into the application bundle.
/// To configure the locales supported by your app, you’ll need to edit this
/// file.
///
/// First, open your project’s ios/Runner.xcworkspace Xcode workspace file.
/// Then, in the Project Navigator, open the Info.plist file under the Runner
/// project’s Runner folder.
///
/// Next, select the Information Property List item, select Add Item from the
/// Editor menu, then select Localizations from the pop-up menu.
///
/// Select and expand the newly-created Localizations item then, for each
/// locale your application supports, add a new item and select the locale
/// you wish to add from the pop-up menu in the Value field. This list should
/// be consistent with the languages listed in the AppLocalizations.supportedLocales
/// property.
abstract class AppLocalizations {
  AppLocalizations(String locale)
    : localeName = intl.Intl.canonicalizedLocale(locale.toString());

  final String localeName;

  static AppLocalizations of(BuildContext context) {
    return Localizations.of<AppLocalizations>(context, AppLocalizations)!;
  }

  static const LocalizationsDelegate<AppLocalizations> delegate =
      _AppLocalizationsDelegate();

  /// A list of this localizations delegate along with the default localizations
  /// delegates.
  ///
  /// Returns a list of localizations delegates containing this delegate along with
  /// GlobalMaterialLocalizations.delegate, GlobalCupertinoLocalizations.delegate,
  /// and GlobalWidgetsLocalizations.delegate.
  ///
  /// Additional delegates can be added by appending to this list in
  /// MaterialApp. This list does not have to be used at all if a custom list
  /// of delegates is preferred or required.
  static const List<LocalizationsDelegate<dynamic>> localizationsDelegates =
      <LocalizationsDelegate<dynamic>>[
        delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
      ];

  /// A list of this localizations delegate's supported locales.
  static const List<Locale> supportedLocales = <Locale>[
    Locale('en'),
    Locale('ja'),
    Locale('zh'),
  ];

  /// No description provided for @navHistory.
  ///
  /// In en, this message translates to:
  /// **'History'**
  String get navHistory;

  /// No description provided for @navDevices.
  ///
  /// In en, this message translates to:
  /// **'Devices'**
  String get navDevices;

  /// No description provided for @navSettings.
  ///
  /// In en, this message translates to:
  /// **'Settings'**
  String get navSettings;

  /// No description provided for @linkConnecting.
  ///
  /// In en, this message translates to:
  /// **'Connecting…'**
  String get linkConnecting;

  /// No description provided for @linkServiceNotRunning.
  ///
  /// In en, this message translates to:
  /// **'Service not running'**
  String get linkServiceNotRunning;

  /// No description provided for @linkNoServer.
  ///
  /// In en, this message translates to:
  /// **'No server connection'**
  String get linkNoServer;

  /// No description provided for @linkOnlinePeers.
  ///
  /// In en, this message translates to:
  /// **'{count, plural, =1{1 device online} other{{count} devices online}}'**
  String linkOnlinePeers(int count);

  /// No description provided for @linkServerConnected.
  ///
  /// In en, this message translates to:
  /// **'Connected to the server'**
  String get linkServerConnected;

  /// No description provided for @linkConnectionLost.
  ///
  /// In en, this message translates to:
  /// **'Connection lost: {detail}'**
  String linkConnectionLost(String detail);

  /// No description provided for @linkServiceDisconnected.
  ///
  /// In en, this message translates to:
  /// **'Service disconnected'**
  String get linkServiceDisconnected;

  /// No description provided for @linkServiceDisabled.
  ///
  /// In en, this message translates to:
  /// **'Background sync is off'**
  String get linkServiceDisabled;

  /// No description provided for @errServiceNotConnected.
  ///
  /// In en, this message translates to:
  /// **'Not connected to the background service'**
  String get errServiceNotConnected;

  /// No description provided for @errUnimplemented.
  ///
  /// In en, this message translates to:
  /// **'This isn\'t available yet'**
  String get errUnimplemented;

  /// No description provided for @errUnavailable.
  ///
  /// In en, this message translates to:
  /// **'Can\'t reach the background service'**
  String get errUnavailable;

  /// No description provided for @errServiceStart.
  ///
  /// In en, this message translates to:
  /// **'Couldn\'t start the background service: {detail}'**
  String errServiceStart(String detail);

  /// No description provided for @errServiceRestart.
  ///
  /// In en, this message translates to:
  /// **'Couldn\'t restart the background service: {detail}'**
  String errServiceRestart(String detail);

  /// No description provided for @errAutostart.
  ///
  /// In en, this message translates to:
  /// **'Couldn\'t set CopySync to start at login: {detail}'**
  String errAutostart(String detail);

  /// No description provided for @thisMacTitle.
  ///
  /// In en, this message translates to:
  /// **'This Mac'**
  String get thisMacTitle;

  /// No description provided for @thisPcTitle.
  ///
  /// In en, this message translates to:
  /// **'This computer'**
  String get thisPcTitle;

  /// No description provided for @thisMacInText.
  ///
  /// In en, this message translates to:
  /// **'this Mac'**
  String get thisMacInText;

  /// No description provided for @thisPcInText.
  ///
  /// In en, this message translates to:
  /// **'this computer'**
  String get thisPcInText;

  /// No description provided for @bothMacs.
  ///
  /// In en, this message translates to:
  /// **'both Macs'**
  String get bothMacs;

  /// No description provided for @bothDevices.
  ///
  /// In en, this message translates to:
  /// **'both devices'**
  String get bothDevices;

  /// No description provided for @thisDeviceShort.
  ///
  /// In en, this message translates to:
  /// **'This device'**
  String get thisDeviceShort;

  /// No description provided for @fromDevice.
  ///
  /// In en, this message translates to:
  /// **'From {name}'**
  String fromDevice(String name);

  /// No description provided for @cancel.
  ///
  /// In en, this message translates to:
  /// **'Cancel'**
  String get cancel;

  /// No description provided for @close.
  ///
  /// In en, this message translates to:
  /// **'Close'**
  String get close;

  /// No description provided for @save.
  ///
  /// In en, this message translates to:
  /// **'Save'**
  String get save;

  /// No description provided for @retry.
  ///
  /// In en, this message translates to:
  /// **'Try again'**
  String get retry;

  /// No description provided for @continueAction.
  ///
  /// In en, this message translates to:
  /// **'Continue'**
  String get continueAction;

  /// No description provided for @justNow.
  ///
  /// In en, this message translates to:
  /// **'Just now'**
  String get justNow;

  /// No description provided for @secondsAgo.
  ///
  /// In en, this message translates to:
  /// **'{n}s ago'**
  String secondsAgo(int n);

  /// No description provided for @minutesAgo.
  ///
  /// In en, this message translates to:
  /// **'{n, plural, =1{1 minute ago} other{{n} minutes ago}}'**
  String minutesAgo(int n);

  /// No description provided for @hoursAgo.
  ///
  /// In en, this message translates to:
  /// **'{n, plural, =1{1 hour ago} other{{n} hours ago}}'**
  String hoursAgo(int n);

  /// No description provided for @daysAgo.
  ///
  /// In en, this message translates to:
  /// **'{n, plural, =1{Yesterday} other{{n} days ago}}'**
  String daysAgo(int n);

  /// No description provided for @monthDay.
  ///
  /// In en, this message translates to:
  /// **'{month}/{day}'**
  String monthDay(int month, int day);

  /// No description provided for @durationHours.
  ///
  /// In en, this message translates to:
  /// **'{n, plural, =1{1 hour} other{{n} hours}}'**
  String durationHours(int n);

  /// No description provided for @durationDays.
  ///
  /// In en, this message translates to:
  /// **'{n, plural, =1{1 day} other{{n} days}}'**
  String durationDays(int n);

  /// No description provided for @historyEmptySubtitle.
  ///
  /// In en, this message translates to:
  /// **'Copy on any of your devices and it shows up here'**
  String get historyEmptySubtitle;

  /// No description provided for @historyCount.
  ///
  /// In en, this message translates to:
  /// **'{count, plural, =1{1 item} other{{count} items}}'**
  String historyCount(int count);

  /// No description provided for @historyCountKept.
  ///
  /// In en, this message translates to:
  /// **'{count, plural, =1{1 item} other{{count} items}}, kept for {ttl}'**
  String historyCountKept(int count, String ttl);

  /// No description provided for @filterAll.
  ///
  /// In en, this message translates to:
  /// **'All'**
  String get filterAll;

  /// No description provided for @filterText.
  ///
  /// In en, this message translates to:
  /// **'Text'**
  String get filterText;

  /// No description provided for @filterImage.
  ///
  /// In en, this message translates to:
  /// **'Images'**
  String get filterImage;

  /// No description provided for @filterFile.
  ///
  /// In en, this message translates to:
  /// **'Files'**
  String get filterFile;

  /// No description provided for @historyClearTooltip.
  ///
  /// In en, this message translates to:
  /// **'Clear history'**
  String get historyClearTooltip;

  /// No description provided for @historyEmptyTitle.
  ///
  /// In en, this message translates to:
  /// **'Nothing here yet'**
  String get historyEmptyTitle;

  /// No description provided for @historyEmptyFilteredTitle.
  ///
  /// In en, this message translates to:
  /// **'Nothing of this kind'**
  String get historyEmptyFilteredTitle;

  /// No description provided for @historyEmptyBody.
  ///
  /// In en, this message translates to:
  /// **'Copy text, an image or a file on any paired device and it appears here.'**
  String get historyEmptyBody;

  /// No description provided for @clearConfirmTitle.
  ///
  /// In en, this message translates to:
  /// **'Clear the whole history?'**
  String get clearConfirmTitle;

  /// No description provided for @clearConfirmBody.
  ///
  /// In en, this message translates to:
  /// **'Everything synced from your devices, and the cached files, will be removed from {device}. The original files aren\'t affected.'**
  String clearConfirmBody(String device);

  /// No description provided for @clear.
  ///
  /// In en, this message translates to:
  /// **'Clear'**
  String get clear;

  /// No description provided for @historyCleared.
  ///
  /// In en, this message translates to:
  /// **'History cleared'**
  String get historyCleared;

  /// No description provided for @permDeniedTitle.
  ///
  /// In en, this message translates to:
  /// **'CopySync can\'t read the clipboard'**
  String get permDeniedTitle;

  /// No description provided for @permAskTitle.
  ///
  /// In en, this message translates to:
  /// **'Allow CopySync to read the clipboard'**
  String get permAskTitle;

  /// No description provided for @permDeniedBody.
  ///
  /// In en, this message translates to:
  /// **'What you copy on this Mac won\'t be synced. In System Settings, set CopySync Daemon to Allow.'**
  String get permDeniedBody;

  /// No description provided for @permAskBody.
  ///
  /// In en, this message translates to:
  /// **'Otherwise macOS asks every time. In System Settings, set CopySync Daemon to Allow.'**
  String get permAskBody;

  /// No description provided for @openSystemSettings.
  ///
  /// In en, this message translates to:
  /// **'Open System Settings'**
  String get openSystemSettings;

  /// No description provided for @fetch.
  ///
  /// In en, this message translates to:
  /// **'Pull to this device'**
  String get fetch;

  /// No description provided for @fetchStarted.
  ///
  /// In en, this message translates to:
  /// **'Pulling…'**
  String get fetchStarted;

  /// No description provided for @transferFailed.
  ///
  /// In en, this message translates to:
  /// **'Transfer failed'**
  String get transferFailed;

  /// No description provided for @expired.
  ///
  /// In en, this message translates to:
  /// **'Expired'**
  String get expired;

  /// No description provided for @preview.
  ///
  /// In en, this message translates to:
  /// **'Preview'**
  String get preview;

  /// No description provided for @putOnClipboard.
  ///
  /// In en, this message translates to:
  /// **'Copy to clipboard'**
  String get putOnClipboard;

  /// No description provided for @putOnClipboardDone.
  ///
  /// In en, this message translates to:
  /// **'On the clipboard'**
  String get putOnClipboardDone;

  /// No description provided for @deleteRecord.
  ///
  /// In en, this message translates to:
  /// **'Delete this item'**
  String get deleteRecord;

  /// No description provided for @deleted.
  ///
  /// In en, this message translates to:
  /// **'Deleted'**
  String get deleted;

  /// No description provided for @emptyRecord.
  ///
  /// In en, this message translates to:
  /// **'(empty)'**
  String get emptyRecord;

  /// No description provided for @imageRecord.
  ///
  /// In en, this message translates to:
  /// **'Image'**
  String get imageRecord;

  /// No description provided for @itemsSummary.
  ///
  /// In en, this message translates to:
  /// **'{name} and {others, plural, =1{1 more} other{{others} more}}'**
  String itemsSummary(String name, int others, int count);

  /// No description provided for @imagePreviewTitle.
  ///
  /// In en, this message translates to:
  /// **'Image preview'**
  String get imagePreviewTitle;

  /// No description provided for @closePreview.
  ///
  /// In en, this message translates to:
  /// **'Close preview'**
  String get closePreview;

  /// No description provided for @previewHint.
  ///
  /// In en, this message translates to:
  /// **'Scroll or pinch to zoom, drag to move around'**
  String get previewHint;

  /// No description provided for @fitWindow.
  ///
  /// In en, this message translates to:
  /// **'Fit to window'**
  String get fitWindow;

  /// No description provided for @previewDeleted.
  ///
  /// In en, this message translates to:
  /// **'This image has been deleted from the history'**
  String get previewDeleted;

  /// No description provided for @previewExpired.
  ///
  /// In en, this message translates to:
  /// **'This image has expired and can\'t be previewed'**
  String get previewExpired;

  /// No description provided for @previewDownloading.
  ///
  /// In en, this message translates to:
  /// **'Downloading the image…'**
  String get previewDownloading;

  /// No description provided for @previewDownloadFailed.
  ///
  /// In en, this message translates to:
  /// **'The download failed. Try again.'**
  String get previewDownloadFailed;

  /// No description provided for @previewNotDownloaded.
  ///
  /// In en, this message translates to:
  /// **'This image isn\'t on this computer yet'**
  String get previewNotDownloaded;

  /// No description provided for @downloadAndPreview.
  ///
  /// In en, this message translates to:
  /// **'Download and preview'**
  String get downloadAndPreview;

  /// No description provided for @previewUnavailable.
  ///
  /// In en, this message translates to:
  /// **'This image can\'t be previewed right now'**
  String get previewUnavailable;

  /// No description provided for @imageSemantic.
  ///
  /// In en, this message translates to:
  /// **'Image from the history'**
  String get imageSemantic;

  /// No description provided for @previewUnreadable.
  ///
  /// In en, this message translates to:
  /// **'Can\'t read the image. The file may have been cleaned up or damaged.'**
  String get previewUnreadable;

  /// No description provided for @devicesSubtitleEmpty.
  ///
  /// In en, this message translates to:
  /// **'Once paired, your devices share one clipboard'**
  String get devicesSubtitleEmpty;

  /// No description provided for @devicesSubtitle.
  ///
  /// In en, this message translates to:
  /// **'{paired} paired, {online} online'**
  String devicesSubtitle(int paired, int online);

  /// No description provided for @addDevice.
  ///
  /// In en, this message translates to:
  /// **'Add device'**
  String get addDevice;

  /// No description provided for @thisDeviceFootnote.
  ///
  /// In en, this message translates to:
  /// **'When pairing, {both} show the same two lines of security fingerprints, one for each device. Confirm only if they match character for character.'**
  String thisDeviceFootnote(String both);

  /// No description provided for @noDevicesTitle.
  ///
  /// In en, this message translates to:
  /// **'No paired devices yet'**
  String get noDevicesTitle;

  /// No description provided for @noDevicesBody.
  ///
  /// In en, this message translates to:
  /// **'Open CopySync on your other computer and link the two with a 6-character pairing code.'**
  String get noDevicesBody;

  /// No description provided for @noDevicesNeedServer.
  ///
  /// In en, this message translates to:
  /// **'Connect to a signaling server in Settings before pairing.'**
  String get noDevicesNeedServer;

  /// No description provided for @pairedDevices.
  ///
  /// In en, this message translates to:
  /// **'Paired devices'**
  String get pairedDevices;

  /// No description provided for @pairedFootnote.
  ///
  /// In en, this message translates to:
  /// **'Direct: data goes straight between the two devices. Relay: when the network blocks a direct path, data goes through your server and stays end-to-end encrypted.'**
  String get pairedFootnote;

  /// No description provided for @offline.
  ///
  /// In en, this message translates to:
  /// **'Offline'**
  String get offline;

  /// No description provided for @direct.
  ///
  /// In en, this message translates to:
  /// **'Direct'**
  String get direct;

  /// No description provided for @relay.
  ///
  /// In en, this message translates to:
  /// **'Relay'**
  String get relay;

  /// No description provided for @online.
  ///
  /// In en, this message translates to:
  /// **'Online'**
  String get online;

  /// No description provided for @moreActions.
  ///
  /// In en, this message translates to:
  /// **'More'**
  String get moreActions;

  /// No description provided for @copyFingerprint.
  ///
  /// In en, this message translates to:
  /// **'Copy fingerprint'**
  String get copyFingerprint;

  /// No description provided for @unpairEllipsis.
  ///
  /// In en, this message translates to:
  /// **'Unpair…'**
  String get unpairEllipsis;

  /// No description provided for @fingerprintCopied.
  ///
  /// In en, this message translates to:
  /// **'Fingerprint copied'**
  String get fingerprintCopied;

  /// No description provided for @unpairTitle.
  ///
  /// In en, this message translates to:
  /// **'Unpair “{name}”?'**
  String unpairTitle(String name);

  /// No description provided for @unpairBody.
  ///
  /// In en, this message translates to:
  /// **'The two devices will stop syncing. To sync again, you\'ll need to pair them and check the fingerprints again.'**
  String get unpairBody;

  /// No description provided for @unpair.
  ///
  /// In en, this message translates to:
  /// **'Unpair'**
  String get unpair;

  /// No description provided for @unpaired.
  ///
  /// In en, this message translates to:
  /// **'Unpaired'**
  String get unpaired;

  /// No description provided for @generateCode.
  ///
  /// In en, this message translates to:
  /// **'Show a code'**
  String get generateCode;

  /// No description provided for @enterCode.
  ///
  /// In en, this message translates to:
  /// **'Enter a code'**
  String get enterCode;

  /// No description provided for @codeFailedTitle.
  ///
  /// In en, this message translates to:
  /// **'Couldn\'t create a pairing code'**
  String get codeFailedTitle;

  /// No description provided for @enterOnOther.
  ///
  /// In en, this message translates to:
  /// **'Enter this code on your other device'**
  String get enterOnOther;

  /// No description provided for @codeCopied.
  ///
  /// In en, this message translates to:
  /// **'Code copied'**
  String get codeCopied;

  /// No description provided for @codeHint.
  ///
  /// In en, this message translates to:
  /// **'Click to copy. Valid for 5 minutes.'**
  String get codeHint;

  /// No description provided for @afterEnterFootnote.
  ///
  /// In en, this message translates to:
  /// **'Once it\'s entered, {both} show the same two lines of security fingerprints. Confirm only if they match.'**
  String afterEnterFootnote(String both);

  /// No description provided for @codeInvalid.
  ///
  /// In en, this message translates to:
  /// **'A pairing code is 6 letters or digits'**
  String get codeInvalid;

  /// No description provided for @enterCodeShown.
  ///
  /// In en, this message translates to:
  /// **'Enter the code shown on your other device'**
  String get enterCodeShown;

  /// No description provided for @pairingDone.
  ///
  /// In en, this message translates to:
  /// **'Paired'**
  String get pairingDone;

  /// No description provided for @pairingRejected.
  ///
  /// In en, this message translates to:
  /// **'Pairing declined'**
  String get pairingRejected;

  /// No description provided for @verifyTitle.
  ///
  /// In en, this message translates to:
  /// **'Check the security fingerprints'**
  String get verifyTitle;

  /// No description provided for @pairingWith.
  ///
  /// In en, this message translates to:
  /// **'Pairing with “{name}”'**
  String pairingWith(String name);

  /// No description provided for @verifyBody.
  ///
  /// In en, this message translates to:
  /// **'The two lines on {both} must be identical. If even one character differs, someone may be tampering with the connection, so decline.'**
  String verifyBody(String both);

  /// No description provided for @rejectMismatch.
  ///
  /// In en, this message translates to:
  /// **'Mismatch, decline'**
  String get rejectMismatch;

  /// No description provided for @confirmMatch.
  ///
  /// In en, this message translates to:
  /// **'They match, pair'**
  String get confirmMatch;

  /// No description provided for @server.
  ///
  /// In en, this message translates to:
  /// **'Server'**
  String get server;

  /// No description provided for @serverFootnote.
  ///
  /// In en, this message translates to:
  /// **'The server only helps devices find each other and set up a direct connection. It can\'t see what they send.'**
  String get serverFootnote;

  /// No description provided for @signalingUrl.
  ///
  /// In en, this message translates to:
  /// **'Signaling server'**
  String get signalingUrl;

  /// No description provided for @signalingUrlDesc.
  ///
  /// In en, this message translates to:
  /// **'Press Return to save; it reconnects right away'**
  String get signalingUrlDesc;

  /// No description provided for @signalingUrlPlaceholder.
  ///
  /// In en, this message translates to:
  /// **'ws://server-address:8787/signal'**
  String get signalingUrlPlaceholder;

  /// No description provided for @publicStun.
  ///
  /// In en, this message translates to:
  /// **'Probe with public servers'**
  String get publicStun;

  /// No description provided for @publicStunDesc.
  ///
  /// In en, this message translates to:
  /// **'On networks with several uplinks, such as offices with two ISPs, this finds more direct paths. Public servers only see your public IP.'**
  String get publicStunDesc;

  /// No description provided for @connectionStatus.
  ///
  /// In en, this message translates to:
  /// **'Connection'**
  String get connectionStatus;

  /// No description provided for @connectedDesc.
  ///
  /// In en, this message translates to:
  /// **'Connected to the signaling server'**
  String get connectedDesc;

  /// No description provided for @disconnectedDesc.
  ///
  /// In en, this message translates to:
  /// **'Can\'t reach the server. Check the address and that the server is running.'**
  String get disconnectedDesc;

  /// No description provided for @connected.
  ///
  /// In en, this message translates to:
  /// **'Connected'**
  String get connected;

  /// No description provided for @notConnected.
  ///
  /// In en, this message translates to:
  /// **'Not connected'**
  String get notConnected;

  /// No description provided for @deviceName.
  ///
  /// In en, this message translates to:
  /// **'Device name'**
  String get deviceName;

  /// No description provided for @deviceNameDesc.
  ///
  /// In en, this message translates to:
  /// **'Shown in the list on your other devices'**
  String get deviceNameDesc;

  /// No description provided for @language.
  ///
  /// In en, this message translates to:
  /// **'Language'**
  String get language;

  /// No description provided for @languageSystem.
  ///
  /// In en, this message translates to:
  /// **'System'**
  String get languageSystem;

  /// No description provided for @sync.
  ///
  /// In en, this message translates to:
  /// **'Sync'**
  String get sync;

  /// No description provided for @syncFootnote.
  ///
  /// In en, this message translates to:
  /// **'Items under the limit are pushed to your other devices as you copy. Larger ones sync as a record; pull them from History when you need them, so big files don\'t waste bandwidth or disk space.'**
  String get syncFootnote;

  /// No description provided for @autoSyncLimit.
  ///
  /// In en, this message translates to:
  /// **'Auto-sync limit'**
  String get autoSyncLimit;

  /// No description provided for @autoSyncLimitDesc.
  ///
  /// In en, this message translates to:
  /// **'Larger files are pulled by hand'**
  String get autoSyncLimitDesc;

  /// No description provided for @autoApply.
  ///
  /// In en, this message translates to:
  /// **'Put received items on the clipboard'**
  String get autoApply;

  /// No description provided for @autoApplyDesc.
  ///
  /// In en, this message translates to:
  /// **'When off, put them on the clipboard from History yourself'**
  String get autoApplyDesc;

  /// No description provided for @whatToSync.
  ///
  /// In en, this message translates to:
  /// **'What to sync'**
  String get whatToSync;

  /// No description provided for @syncFiles.
  ///
  /// In en, this message translates to:
  /// **'Files and folders'**
  String get syncFiles;

  /// No description provided for @syncText.
  ///
  /// In en, this message translates to:
  /// **'Plain text'**
  String get syncText;

  /// No description provided for @syncRich.
  ///
  /// In en, this message translates to:
  /// **'Rich text'**
  String get syncRich;

  /// No description provided for @syncImages.
  ///
  /// In en, this message translates to:
  /// **'Images and screenshots'**
  String get syncImages;

  /// No description provided for @storage.
  ///
  /// In en, this message translates to:
  /// **'Storage'**
  String get storage;

  /// No description provided for @storageFootnote.
  ///
  /// In en, this message translates to:
  /// **'Expired items and cached files are cleaned up automatically.'**
  String get storageFootnote;

  /// No description provided for @keepHistory.
  ///
  /// In en, this message translates to:
  /// **'Keep history for'**
  String get keepHistory;

  /// No description provided for @keepCache.
  ///
  /// In en, this message translates to:
  /// **'Keep cached files for'**
  String get keepCache;

  /// No description provided for @keepCacheDesc.
  ///
  /// In en, this message translates to:
  /// **'Never longer than the history'**
  String get keepCacheDesc;

  /// No description provided for @cacheUsage.
  ///
  /// In en, this message translates to:
  /// **'Cache size'**
  String get cacheUsage;

  /// No description provided for @backgroundSync.
  ///
  /// In en, this message translates to:
  /// **'Background sync'**
  String get backgroundSync;

  /// No description provided for @backgroundFootnote.
  ///
  /// In en, this message translates to:
  /// **'When off, nothing syncs and CopySync doesn\'t start at login. History and pairings are kept; open the app to turn it back on.'**
  String get backgroundFootnote;

  /// No description provided for @backgroundService.
  ///
  /// In en, this message translates to:
  /// **'Background service'**
  String get backgroundService;

  /// No description provided for @backgroundServiceDesc.
  ///
  /// In en, this message translates to:
  /// **'Starts at login and keeps syncing with the window closed'**
  String get backgroundServiceDesc;

  /// No description provided for @restart.
  ///
  /// In en, this message translates to:
  /// **'Restart'**
  String get restart;

  /// No description provided for @restarted.
  ///
  /// In en, this message translates to:
  /// **'Background service restarted'**
  String get restarted;

  /// No description provided for @disableEllipsis.
  ///
  /// In en, this message translates to:
  /// **'Turn off…'**
  String get disableEllipsis;

  /// No description provided for @about.
  ///
  /// In en, this message translates to:
  /// **'About'**
  String get about;

  /// No description provided for @aboutDesc.
  ///
  /// In en, this message translates to:
  /// **'Open source under the AGPL-3.0'**
  String get aboutDesc;

  /// No description provided for @homepage.
  ///
  /// In en, this message translates to:
  /// **'Website'**
  String get homepage;

  /// No description provided for @feedback.
  ///
  /// In en, this message translates to:
  /// **'Report a problem'**
  String get feedback;

  /// No description provided for @disableTitle.
  ///
  /// In en, this message translates to:
  /// **'Turn off background sync?'**
  String get disableTitle;

  /// No description provided for @disableBody.
  ///
  /// In en, this message translates to:
  /// **'CopySync on {device} will stop syncing and won\'t start at login. History and pairings are kept.'**
  String disableBody(String device);

  /// No description provided for @disable.
  ///
  /// In en, this message translates to:
  /// **'Turn off'**
  String get disable;

  /// No description provided for @disabled.
  ///
  /// In en, this message translates to:
  /// **'Background sync is off'**
  String get disabled;

  /// No description provided for @urlMustBeWs.
  ///
  /// In en, this message translates to:
  /// **'The address must start with ws:// or wss://'**
  String get urlMustBeWs;

  /// No description provided for @serverUpdated.
  ///
  /// In en, this message translates to:
  /// **'Server updated. Reconnecting…'**
  String get serverUpdated;

  /// No description provided for @deviceNameUpdated.
  ///
  /// In en, this message translates to:
  /// **'Device name updated'**
  String get deviceNameUpdated;

  /// No description provided for @svcExtractTitle.
  ///
  /// In en, this message translates to:
  /// **'Install or extract CopySync first'**
  String get svcExtractTitle;

  /// No description provided for @svcExtractBody.
  ///
  /// In en, this message translates to:
  /// **'CopySync is running straight from the zip, and its temporary folder disappears when you close it. Run the installer, or extract the whole folder and open it from there, to turn on background sync.'**
  String get svcExtractBody;

  /// No description provided for @svcMoveTitle.
  ///
  /// In en, this message translates to:
  /// **'Move CopySync to Applications first'**
  String get svcMoveTitle;

  /// No description provided for @svcMoveBody.
  ///
  /// In en, this message translates to:
  /// **'CopySync is running straight from the disk image. Drag it into Applications, eject the disk image, and open it from Applications to turn on background sync.'**
  String get svcMoveBody;

  /// No description provided for @svcStartTitle.
  ///
  /// In en, this message translates to:
  /// **'Get started with CopySync'**
  String get svcStartTitle;

  /// No description provided for @svcStartBody.
  ///
  /// In en, this message translates to:
  /// **'CopySync runs a small background process that notices clipboard changes and stays connected to your other devices. Once on, it starts at login and keeps syncing even with this window closed.'**
  String get svcStartBody;

  /// No description provided for @svcEnable.
  ///
  /// In en, this message translates to:
  /// **'Turn on background sync'**
  String get svcEnable;

  /// No description provided for @svcStaleTitle.
  ///
  /// In en, this message translates to:
  /// **'Turn background sync on again'**
  String get svcStaleTitle;

  /// No description provided for @svcStaleBodyWin.
  ///
  /// In en, this message translates to:
  /// **'CopySync has moved, but starting at login still points to the old place. Turn it on again so the background process uses this copy. History and pairings are kept.'**
  String get svcStaleBodyWin;

  /// No description provided for @svcStaleBodyMac.
  ///
  /// In en, this message translates to:
  /// **'CopySync has moved, or was installed with an old script. Turn it on again so the background process uses this app. History and pairings are kept.'**
  String get svcStaleBodyMac;

  /// No description provided for @svcReenable.
  ///
  /// In en, this message translates to:
  /// **'Turn on again'**
  String get svcReenable;

  /// No description provided for @svcStoppedTitle.
  ///
  /// In en, this message translates to:
  /// **'The background service isn\'t running'**
  String get svcStoppedTitle;

  /// No description provided for @svcStoppedBody.
  ///
  /// In en, this message translates to:
  /// **'It may have just quit; the system usually restarts it within seconds. If you stay on this page, restart it with the button below.'**
  String get svcStoppedBody;

  /// No description provided for @svcDevBodyWin.
  ///
  /// In en, this message translates to:
  /// **'This is a development build with no copysyncd.exe next to it. Run scripts\\build-windows.ps1 in the source folder to build a complete version.'**
  String get svcDevBodyWin;

  /// No description provided for @svcDevBodyMac.
  ///
  /// In en, this message translates to:
  /// **'This is a development build without the background service. Run ./scripts/install-macos.sh in the source folder to install it.'**
  String get svcDevBodyMac;

  /// No description provided for @svcLogPath.
  ///
  /// In en, this message translates to:
  /// **'Log: {path}'**
  String svcLogPath(String path);

  /// No description provided for @svcCanDisable.
  ///
  /// In en, this message translates to:
  /// **'You can turn background sync off anytime in Settings.'**
  String get svcCanDisable;
}

class _AppLocalizationsDelegate
    extends LocalizationsDelegate<AppLocalizations> {
  const _AppLocalizationsDelegate();

  @override
  Future<AppLocalizations> load(Locale locale) {
    return SynchronousFuture<AppLocalizations>(lookupAppLocalizations(locale));
  }

  @override
  bool isSupported(Locale locale) =>
      <String>['en', 'ja', 'zh'].contains(locale.languageCode);

  @override
  bool shouldReload(_AppLocalizationsDelegate old) => false;
}

AppLocalizations lookupAppLocalizations(Locale locale) {
  // Lookup logic when only language code is specified.
  switch (locale.languageCode) {
    case 'en':
      return AppLocalizationsEn();
    case 'ja':
      return AppLocalizationsJa();
    case 'zh':
      return AppLocalizationsZh();
  }

  throw FlutterError(
    'AppLocalizations.delegate failed to load unsupported locale "$locale". This is likely '
    'an issue with the localizations generation tool. Please file an issue '
    'on GitHub with a reproducible sample app and the gen-l10n configuration '
    'that was used.',
  );
}
