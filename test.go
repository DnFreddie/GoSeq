package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
)

// bufferHash calculates SHA-256 hash of buffer contents
func bufferHash(buf *bytes.Buffer) string {
	hasher := sha256.New()
	hasher.Write(buf.Bytes())
	return hex.EncodeToString(hasher.Sum(nil))
}
