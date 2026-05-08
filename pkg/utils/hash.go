package utils

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
)

func MD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

func HashPassword(password, salt string) string {
	passwordBytes := append([]byte(password), salt...)

	h := sha512.New()
	h.Write(passwordBytes)
	hashedBytes := h.Sum(nil)
	return hex.EncodeToString(hashedBytes)
}
