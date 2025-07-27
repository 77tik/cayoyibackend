package util

import "crypto/rand"

func RandomBytes(byteCount int) []byte {
	buf := make([]byte, byteCount)
	rand.Read(buf)
	return buf
}
