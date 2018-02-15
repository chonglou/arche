package nut

import (
	"crypto/sha512"
	"encoding/base64"
)

// SumSsha512 sum ssha512
func SumSsha512(plain string, sl int) (string, error) {
	salt, err := RandomBytes(sl)
	if err != nil {
		return "", err
	}
	return sumSsha512(plain, salt), nil
}

// EqualSsha512 compare ssha512
func EqualSsha512(encode, plain string) bool {
	if buf, err := base64.StdEncoding.DecodeString(encode); err == nil {
		return encode == sumSsha512(plain, buf[sha512.Size:])
	}
	return false
}

func sumSsha512(plain string, salt []byte) string {
	buf := sha512.Sum512(append([]byte(plain), salt...))
	return base64.StdEncoding.EncodeToString(append(buf[:], salt...))
}
