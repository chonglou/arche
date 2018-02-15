package nut

import (
	"encoding/base64"
	"net/http"
	"sync"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/SermoDigital/jose/jwt"
	"github.com/astaxie/beego"
)

var (
	_jwt     *Jwt
	_jwtOnce sync.Once
)

// JWT jwt instance
func JWT() *Jwt {
	_jwtOnce.Do(func() {
		key, err := base64.StdEncoding.DecodeString(beego.AppConfig.String("secrets"))
		if err != nil {
			beego.Error(err)
			return
		}
		_jwt = &Jwt{key: key, method: crypto.SigningMethodHS512}
	})
	return _jwt
}

// Jwt jwt
type Jwt struct {
	key    []byte
	method crypto.SigningMethod
}

//Validate check jwt
func (p *Jwt) Validate(buf []byte) (jwt.Claims, error) {
	tk, err := jws.ParseJWT(buf)
	if err != nil {
		return nil, err
	}
	if err = tk.Validate(p.key, p.method); err != nil {
		return nil, err
	}
	return tk.Claims(), nil
}

// Parse parse
func (p *Jwt) Parse(r *http.Request) (jwt.Claims, error) {
	tk, err := jws.ParseJWTFromRequest(r)
	if err != nil {
		return nil, err
	}
	if err = tk.Validate(p.key, p.method); err != nil {
		return nil, err
	}
	return tk.Claims(), nil
}

//Sum create jwt token
func (p *Jwt) Sum(cm jws.Claims, exp time.Duration) ([]byte, error) {
	now := time.Now()
	cm.SetNotBefore(now)
	cm.SetExpiration(now.Add(exp))

	jt := jws.NewJWT(cm, p.method)
	return jt.Serialize(p.key)
}
