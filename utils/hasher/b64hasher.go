package b64hasher

import (
	"crypto/sha512"
	"encoding/base64"
	"math/rand"
)

func GenerateSalt() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+="
	salt := make([]byte, 32)
	for i := range salt {
		salt[i] = charset[rand.Intn(len(charset))]
	}
	return string(salt)
}

func HashPassword(password, salt string) string {
	passwordBytes := []byte(password + salt)

	sha512Hasher := sha512.New()
	sha512Hasher.Write(passwordBytes)

	hashedPasswordBytes := sha512Hasher.Sum(nil)
	b64encoded := b64Encode(hashedPasswordBytes)
	return b64encoded
}

func IsPasswordsMatch(password, salt, hash string) bool {
	currPasswordHash := HashPassword(password, salt)
	return b64Decode(hash) == b64Decode(currPasswordHash)
}

func IsPasswordHashed(password string) bool {
	decoded := b64Decode(password)
	if len(decoded) == 0 {
		return false
	}
	return true
}

func b64Encode(encoded []byte) string {
	cstr := base64.URLEncoding.EncodeToString(encoded)
	return cstr
}

func b64Decode(encoded string) string {
	cstr, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return ""
	}
	return string(cstr)
}
