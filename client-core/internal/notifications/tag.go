package notifications

import (
	"crypto/sha256"
	"encoding/hex"
)

func notificationTag(id string) string {
	digest := sha256.Sum256([]byte(id))
	return hex.EncodeToString(digest[:8])
}
