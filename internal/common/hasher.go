package common

import (
	"crypto/sha256"
	"encoding/hex"
)


func calculateHash(content []byte) string {
	hasher := sha256.New()
	hasher.Write(content)
	return hex.EncodeToString(hasher.Sum(nil))
}
