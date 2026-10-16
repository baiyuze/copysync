# Security model and reporting

PasteFerry clients and the signaling/TURN server are open source under AGPL-3.0.
Self-hosting never requires an official account, authorization service or
closed-source component. See [the optional gateway protocol](design/official-gateway.md).

## Trust boundaries

- Device identity is an existing local random Ed25519 key. Pairing requires
  comparing and confirming fingerprints. Server announcements alone do not
  establish trust. A machine quota digest is neither a credential nor a key seed.
- Device-to-device SDP fingerprints and signaling are signed. DTLS protects
  clipboard/file data in transit, including over relay. The gateway can see
  transport metadata and byte counts; it cannot read those encrypted contents.
- Official access is optional and disabled for fresh installs. LAN and custom
  modes do not request official authorization or report official usage.
- The gateway can allow, deny, rate-limit or interrupt its own hosted network
  resources. It cannot execute client code, read arbitrary local files, export
  device keys, silently pair devices or change the cryptographic trust model.
- Protocol parsing, body/frame limits, expiration, replay checks and destination
  restrictions must remain client- and server-side validation, regardless of
  the operator's implementation. Official endpoints are public configuration.
  No server credential belongs in an application bundle.

## Practical limits

An administrator or modified client can replace machine identifiers. Hashed
device metadata is linkable and is not an immutable hardware identity or proof
of a unique person. Timing, IP addresses and ICE routing metadata are visible
to the relevant network endpoints during connection establishment. LAN
discovery intentionally advertises local presence. Physical compromise of a
device can expose its local clipboard history and cached files.

Ordinary Cloudflare HTTP proxying does not protect UDP STUN or carry TURN UDP.
Quota limits control authorized service use, not provider-level DDoS ingress.
Source review establishes the published application's behavior; it cannot by
itself prove a private service's actual production retention or deployment.

## Verify a build

Review the release commit, dependency lock files, CI results and SHA256SUMS.
Build instructions are in `scripts/build.sh`, `scripts/package.sh` and the
server deployment documentation. A checksum confirms artifact integrity; it
does not prove reproducible compilation. Do not claim reproducibility unless
the actual platform artifacts have been independently compared.

Recommended checks: `go vet ./...`, `go test -race ./...` in Go modules,
`govulncheck ./...`, `flutter analyze`, `flutter test`, and the black-box
gateway interoperability tests when an independently built gateway is supplied.

## Report a vulnerability

Use [GitHub private vulnerability reporting](https://github.com/baiyuze/pasteferry/security/advisories/new).
Describe affected versions, reproduction steps and impact. Do not put device
private keys, live authorization tokens, passwords or clipboard contents in a
public issue. If private reporting is unavailable, open an issue requesting a
private contact without publishing exploitation details.
