package utils

import (
	"bytes"
	mathrand "math/rand"
	"time"

	cryptorand "crypto/rand"
	"encoding/base64"
)

func GenerateMathRandString(size int64) string {
	src := mathrand.NewSource(time.Now().UnixNano())
	r := mathrand.New(src)

	strBytes := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := []byte{}
	for range size {
		idx := r.Intn(bytes.Count(strBytes, nil) - 1)
		temp := strBytes[idx]
		result = append(result, temp)
	}
	return string(result)
}

// Generate 16 bytes randomly and securely
// using the Cryptographically secure pseudorandom number generator (CSPRNG)
// in the crypto.rand package
func GenerateCryptoRandString(size int) string {
	var bytes = make([]byte, size)
	cryptorand.Read(bytes[:])

	return base64.RawStdEncoding.EncodeToString(bytes)
}

func GenerateCryptoRandNumString(n int) string {
	var bytes = make([]byte, n)
	cryptorand.Read(bytes)

	var numstr = "0123456789"
	lenn := len(numstr)

	for i, v := range bytes {
		idx := v % byte(lenn)
		bytes[i] = numstr[idx]
	}
	return string(bytes)
}
