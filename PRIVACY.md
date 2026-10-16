# Privacy and official service data

## Local and self-hosted use

New installations start with official public service disabled. Clipboard
history and cached files are stored locally. LAN discovery and synchronization
use the local network. Custom servers are user-selected; their operators have
their own retention practices. Neither mode uploads official usage telemetry.
Optional public STUN probing reveals the source network address to selected
STUN providers; the setting to use only the chosen server disables extra probes.

## Official access

Enabling the official service is a deliberate setting with an explanation.
The client first uses signaling and STUN to try a direct connection. A failed
direct attempt with an online compatible peer is required before requesting
official relay access. No machine fingerprint is needed for direct service.

| Data | Purpose / retention |
| --- | --- |
| Device public key and signature | Prove possession of the existing identity; gateway authentication state is transient |
| Server-derived pseudonymous device number | Access counts, first/last seen, bans and cumulative totals |
| Application-specific machine digest | Relay quota grouping only; raw system identifier never leaves the client |
| Daily upload/download and access totals | Per-device detail retained for 90 days; essential aggregate totals retained longer |
| Current peer/session and routing metadata | Transient forwarding and abuse protection; not historical relationship tracking |
| IP addresses | Necessarily visible for transport; not persisted in official statistics or ordinary proxy logs |

Real computer names and OS details are not included in new official client
handshakes/pairing requests. Paired devices exchange display names inside the
encrypted data channel. Clipboard content, file names, paths, local usernames,
MAC addresses, browser UA, installed software and behavior analytics are not
fields of the gateway protocol. No third-party client analytics SDK is added.

The machine digest is generated only when relay authorization is requested.
Windows uses MachineGuid; macOS uses IOPlatformUUID; Linux uses a stable
machine-id, never reboot-specific boot_id. A bounded, cached read is wrapped
by a versioned PasteFerry-specific HMAC; the service hashes again with its own
secret. These are pseudonymous identifiers, not guaranteed anonymity. Ordinary
application reinstall generally preserves the quota group; OS reinstall,
cloning or a modified client can change or forge it.

The dashboard counts **official service devices**, not people, total installs,
or users who choose LAN/custom servers. Recreated device keys may count again.

## Usage and controls

The default relay allowance is 100 MB/day, decimal bytes, upload plus download,
reset at midnight Beijing time. Direct transfers do not consume this quota.
Encrypted transport overhead contributes to measured relay data. The gateway
may conservatively hold a small unresolved reservation after a crash until the
daily reset. Unknown/stale usage is labeled unavailable, not displayed as zero.

Turning off official service cancels official requests, polling and tunnels.
There is no separate background analytics channel. Quota exhaustion affects
official relay only; it does not turn off LAN or direct synchronization.
No anonymous-account system can guarantee a unique natural person without
collecting additional identity information; PasteFerry does not attempt that.

## Website and audit

The website has an existing Cloudflare Web Analytics beacon, separate from
the desktop client. Website visits reach Cloudflare and the hosting endpoint;
they are not software installation or active-device counts. Consult the
website provider's privacy policy for its network processing.

The published client, network field list and [gateway protocol](design/official-gateway.md)
can be audited. A private gateway's production data handling requires operator
and deployment audit; open client source alone cannot prove it. See
[SECURITY.md](SECURITY.md) for the trust model and reporting channel.
