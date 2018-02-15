package test

import (
	"testing"

	"github.com/chonglou/arche/plugins/nut"
)

// TestSsha512 test dovecot SSHA512
func TestSsha512(t *testing.T) {
	const plain = "Hi, arche."
	encode, err := nut.SumSsha512(plain, 32)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(`doveadm pw -t {SSHA512}%s -p "%s"`, encode, plain)
	if !nut.EqualSsha512(encode, plain) {
		t.Error("check password failed")
	}
}
