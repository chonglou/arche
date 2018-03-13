package web_test

import (
	"encoding/base64"
	"testing"

	"github.com/chonglou/arche/web"
)

const hello = "Hello, Arche!"

func TestSecurity(t *testing.T) {
	pwd, err := web.RandomBytes(32)
	if err != nil {
		t.Fatal(err)
	}
	sec, err := web.NewSecurity(pwd)
	if err != nil {
		t.Fatal(err)
	}

	testSecurityPassword(t, sec)
	testSecurityEncrypt(t, sec)
}

func testSecurityEncrypt(t *testing.T, s *web.Security) {
	plain := []byte(hello)
	for i := 0; i < 5; i++ {
		encode, err := s.Encrypt(plain)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(base64.StdEncoding.EncodeToString(encode))
		buf, err := s.Decrypt(encode)
		if err != nil {
			t.Fatal(err)
		}
		if string(buf) != hello {
			t.Fatalf("want %s, get %s", hello, string(buf))
		}
	}
}
func testSecurityPassword(t *testing.T, s *web.Security) {
	plain := []byte(hello)
	for i := 0; i < 5; i++ {
		passwd, err := s.Hash(plain)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(passwd))
		if !s.Check(passwd, plain) {
			t.Fatal("test password failed")
		}
	}
}
