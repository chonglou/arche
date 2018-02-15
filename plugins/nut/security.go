package nut

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"sync"

	"github.com/astaxie/beego"
)

var (
	_aes     *Aes
	_aesOnce sync.Once
)

// AES get aes instance.
func AES() *Aes {
	_aesOnce.Do(func() {
		key, err := base64.StdEncoding.DecodeString(beego.AppConfig.String("secrets"))
		if err != nil {
			beego.Error(err)
			return
		}
		cip, err := aes.NewCipher(key)
		if err != nil {
			beego.Error(err)
			return
		}
		_aes = &Aes{cip: cip}
	})

	return _aes
}

// Aes aes helper
type Aes struct {
	cip cipher.Block
}

// Encrypt aes encrypt
func (p *Aes) Encrypt(buf []byte) ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(p.cip, iv)
	val := make([]byte, len(buf))
	cfb.XORKeyStream(val, buf)

	return append(val, iv...), nil
}

// Decrypt aes decrypt
func (p *Aes) Decrypt(buf []byte) ([]byte, error) {
	bln := len(buf)
	cln := bln - aes.BlockSize
	ct := buf[0:cln]
	iv := buf[cln:bln]

	cfb := cipher.NewCFBDecrypter(p.cip, iv)
	val := make([]byte, cln)
	cfb.XORKeyStream(val, ct)
	return val, nil
}

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
