# Isolated cross-platform transport validation

`client-core/cmd/envcheck` exchanges synthetic fixtures between two explicitly
selected endpoints. It does not read the OS clipboard, reuse installed device
keys, change application settings, or inspect existing history. Output must be
a new directory. Test identities are deterministic and publicly reproducible;
never trust or reuse them outside an isolated validation run.

Build on each machine, or cross-compile with `CGO_ENABLED=0` for Windows/Linux:

```sh
cd client-core
go build -o envcheck ./cmd/envcheck
```

Run one process on each endpoint with the same unique run label and reversed
`-self` / `-peer` labels. Both processes must be online together:

```sh
./envcheck -self mac -peer linux -run UNIQUE-TEST-LABEL \
  -mode direct -signal wss://YOUR-TEST-HOST/signal \
  -gateway https://YOUR-TEST-GATEWAY -output NEW-OUTPUT-DIRECTORY
```

Modes:

- `direct`: official signaling with direct ICE candidates. Fails if relay
  authorization is requested. No machine digest is needed.
- `official-relay`: substitutes unusable direct candidates before signing them,
  waits for the real direct-first fallback, then uses the standard TURN service
  through the gateway WSS datagram transport. Fails if the direct attempt window
  is skipped.
- `custom-relay`: forces ordinary TURN with no gateway client. Use an independent
  self-hosted signaling/TURN endpoint; `-gateway` is unused.

Each direction sends multilingual text, a generated PNG and a deterministic
1 MiB file, waits for peer checksum acknowledgements, and writes `result.json`.
Unicode and emoji filenames are saved in the test directory. The default
timeout is three minutes; `-timeout` can adjust a weak-network run.

For an isolated allowance-denial test, `-expect-denial quota_exhausted` succeeds
only when that endpoint observes the specified gateway error. Test the
deterministic offerer: if its authorization is denied before a relay offer, the
other endpoint may never request authorization. Stop that waiting counterpart
after collecting the target result. Never change production user allowances
for this test.

This checks real signaling, signed negotiation, DTLS/DataChannel transport and
file integrity across OS processes. It does not replace native OS clipboard,
installer, GUI interaction or Internet path performance testing.
