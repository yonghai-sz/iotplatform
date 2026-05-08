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
	var salt = make([]byte, size)
	_, err := cryptorand.Read(salt[:])
	if err != nil {
		panic(err)
	}
	return base64.RawStdEncoding.EncodeToString(salt)
}
