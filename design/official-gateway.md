# Optional official gateway protocol v1

The application and the signaling/TURN server remain independent AGPL software.
This document describes an optional, independently implemented service at the
official hosting boundary. Neither server calls the gateway, sends accounting
callbacks, imports a private SDK, or requires the gateway for self-hosting.
The gateway speaks ordinary HTTP/WebSocket and standard TURN on its data path.

## Selection and connection order

Fresh installations use LAN mode. Config v3 adds `service_mode` (`lan`,
`official`, `custom`), a remembered `custom_signaling_url` and local-only notice
acknowledgements. Selecting LAN or custom must cancel official HTTP requests,
usage polling and tunnels. No automatic fallback between service providers.

Official signaling: `wss://rtccpwsapi1.pasteferry.cc/signal`.
Official gateway base: `https://rtccpwsapi1.pasteferry.cc`.
Official STUN is the separate Binding-only endpoint at `104.223.127.214:3478`.
Addresses are public configuration, not credentials. Cloudflare proxies WSS,
not this UDP STUN endpoint. Public STUN probes are separately user-controlled;
their provider addresses are listed in `client-core/internal/transport/p2p/egress.go`.

Initially filter TURN out of the active ICE configuration, retaining it only
in memory. After an online, capable peer exchanges offer/answer and direct
connection fails for eight seconds, the lexically lower device ID creates a
random 32-byte attempt ID. Both devices independently authorize that same
attempt. Only then create relay-only PeerConnections with the WSS packet
adapter. Keep the direct connection independent so a later successful direct
attempt can release the relay. An offline/legacy peer or signaling error is
not evidence that relay authorization should be requested.

Signed SignalMessage adds optional `connection_attempt_id`, `official_relay`
and `gateway_capable` and `relay_denial_code` fields. The server forwards the original signed payload;
it does not interpret these capabilities. PeerDisplay travels only inside the
authenticated encrypted device channel. Official ClientHello and pairing
requests use a generic device label, not the real hostname or OS details.

## Canonical signed proof

All strings are UTF-8 without embedded newlines. Public keys and signatures
use unpadded base64url. Device IDs retain the existing SHA-256(public key),
first 16 bytes, lowercase unpadded base32 derivation. The signed bytes are
these eight values, joined with a single LF, with no final LF:

```
pasteferry-gateway-v1
purpose
nonce
device_id
public_key
peer_id
session_id
machine_digest
```

Sign with the existing device Ed25519 key. This never changes pairing trust.
For signaling, purpose is `signal`, nonce is `unix_seconds.random_base64url`
(24 random bytes), and the final three fields are empty. Accept timestamps
within 60 seconds and reject replay. Send the following headers only to the
official signal endpoint: `X-PasteFerry-Key`, `X-PasteFerry-Device`,
`X-PasteFerry-Nonce`, `X-PasteFerry-Signature`. The proxy removes these headers
before forwarding to the unchanged signaling server. This is an access identity
proof, not a relay quota check. Legacy connections are counted separately.

## Relay and usage API

Every route is under `/gateway/v1/`. HTTPS is mandatory except explicit
loopback interoperability tests. Redirects are rejected. JSON requests are
limited to 8 KiB; response bodies to the client are limited to 32 KiB.

| Method / route | Request | Response |
| --- | --- | --- |
| GET `challenge` | No machine or device information | `nonce`, `expires_in:60` |
| POST `authorize` | Proof fields above plus `signature`; purpose `relay` | `token`, `expires_at`, `usage` |
| GET `usage` | `Authorization: Bearer token` | Authoritative usage snapshot |
| GET `tunnel` | WebSocket upgrade and same Authorization header | Packet connection |

Challenges are random, single-use, and consumed even on invalid signatures.
Session IDs are 32 random bytes in unpadded base64url. Enrollment binds the
current device to its quota group; changing a later digest cannot reset that
binding. Missing machine information retains the existing binding, or falls
back to persistent device identity on first enrollment. Tokens are random,
memory-only, expire after 15 minutes, and bind device, peer, session and group.
Both reciprocal grants must exist before opening a tunnel.

After first relay enrollment, `authorize` also accepts purpose `usage` with
empty peer, session and machine-digest fields and a fresh challenge signature.
It can only issue a read-only usage token for an already-bound device. That
token cannot open or renew a tunnel, and does not create a peer session.
A denied relay grant may include `usage_token` and `expires_at` so the client
can observe reset or administrator changes without retrying relay constantly.
Disabling the official mode clears all tokens and cancels these requests.

A signed `relay_denial_code` is limited to quota/policy denial and the
`quota_available` retry hint. It cannot grant access: the gateway still checks
both devices. It never carries an executable instruction or pairing approval.

Usage fields: `available`, `day`, `group`, `upload`, `download`, `limit`,
`remaining`, `reserved`, `reset_at`, `updated_at`, `reason`. Bytes are integer
decimal bytes; timestamps are Unix seconds. Unknown/stale data is not zero.
Reserved bytes left after a crash remain conservatively unavailable that day.
The server never accepts client-supplied usage totals or a client reset clock.

Failure responses have `code` and optionally `usage`. Stable reasons include
`unauthorized`, `challenge_expired`, `invalid_proof`, `session_mismatch`,
`peer_pending`, `quota_exhausted`, `global_exhausted`, `global_paused`,
`relay_blocked`, `device_banned`, `rate_limited`, `service_unavailable`.
The client maps allowed codes to local translated messages; it does not
execute or directly render remotely supplied instructions.

## Packet tunnel

Each WSS connection represents exactly one connected UDP socket to the fixed
private TURN control address. The first server text frame is JSON containing
`local_address` and `turn_address`; the latter must match the selected TURN
configuration. Binary frames contain one unchanged UDP datagram, up to 65536
bytes. There is no arbitrary-destination field or SOCKS/CONNECT capability.

A client text frame `{ "token": "new token" }` renews authorization on the
same tunnel; it must bind the same device, group, peer and session. Credentials
never appear in query parameters. The tunnel is closed when authorization
expires or gateway policy denies forwarding. Queue lengths and I/O deadlines
are bounded; closing the transport cancels outstanding reads and writes.

The independent proxy observes standard TURN allocation responses and outer
Send/CreatePermission/ChannelBind messages. It permits peer addresses only
when the other allocation belongs to the reciprocally authorized session.
It never decrypts DTLS application records. The TURN server still validates
its own standard credentials and message integrity, without gateway secrets.

## Accounting and failure behavior

Policy inputs use a number with an MB/GB selector, with a minimum adjustment
of 1 MB; speed uses Mbps. 1 GB = 1000 MB. API/storage retain integer base units
so changing display units never changes the actual allowance. Usage is shown
in whole MB (`<1 MB` for positive amounts below 1 MB); packets are accumulated
at their real length, never rounded up individually to 1 MB.

Default quota: 100 MB per machine group per Beijing calendar day,
upload plus download. Shared rate: 2,000,000 bits/s. Each device boundary is
charged once; same-group endpoints share the total. Server-internal TURN hops
are not extra charges. Bill forwarded encrypted payload bytes, including
encryption overhead; file size and usage therefore differ. Normal bounded
TURN/STUN/ICE control exchanges are not charged to the device data quota.

Global budget: 1,000,000,000 outbound tunnel datagram bytes/day, including
control. This measures the gateway application boundary, not TCP/TLS headers,
TCP retransmissions or the hosting provider's bill. Per-group in-memory
enforcement uses bounded durable reservations and idempotent settlement;
never perform a SQLite transaction per packet. Exhaustion blocks existing
and new relay data. Policy changes apply within five seconds.

Auth/ledger outages fail new relay credit closed. Already reserved credit is
the maximum that can be spent before service recovers. A proxy outage affects
official signaling entry, but cannot revoke LAN/direct cryptographic trust or
disable self-hosting. There are no remote-code, file-read, key-export or update
commands in this protocol.

## Deployment contract

Only the private gateway and Nginx are public entry points. The unchanged
signaling/TURN binary runs in a network namespace, advertising its private
relay address. All relay ports are unreachable directly from the Internet;
both clients reach them through TURN encapsulation in the authenticated
tunnels. Backend egress is isolated so it cannot serve as a general proxy.
A separate public Binding-only STUN responder does not forward Allocate.

Custom/self-hosted server binaries need none of this configuration. The
gateway implementation and operator dashboard may be private; this protocol,
client implementation and complete self-hosted server remain public.
