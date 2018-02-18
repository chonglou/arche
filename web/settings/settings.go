package settings

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/go-pg/pg"
)

func New(sec string) (*Settings, error) {
	key, err := base64.StdEncoding.DecodeString(sec)
	if err != nil {
		return nil, err
	}
	cip, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &Settings{cip: cip}, nil
}

// Model locale database model
type Model struct {
	tableName struct{} `sql:"locales"`
	Key       string
	Value     []byte
	Encode    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Settings settings
type Settings struct {
	cip cipher.Block
}

// Set set value
func (p *Settings) Set(db *pg.Tx, k string, v interface{}, f bool) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if f {
		buf, err = p.encrypt(buf)
		if err != nil {
			return err
		}
	}

	var it Model
	err = db.Model(&it).Column("id").Where("key = ?", k).Select()
	it.Value = buf
	it.Encode = f
	it.UpdatedAt = time.Now()
	if err == nil {
		_, err = db.Model(&it).Column("value", "encode", "updated_at").Update()
	} else if err == pg.ErrNoRows {
		it.Key = k
		err = db.Insert(&it)
	}
	return err
}

// Get get value
func (p *Settings) Get(db *pg.DB, k string, v interface{}) error {
	var it Model
	err := db.Model(&it).Column("value", "encode").
		Where("key = ?", k).
		Select()
	if err != nil {
		return err
	}
	if it.Encode {
		it.Value, err = p.decrypt(it.Value)
		if err != nil {
			return err
		}
	}
	return json.Unmarshal(it.Value, v)
}

func (p *Settings) encrypt(buf []byte) ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(p.cip, iv)
	val := make([]byte, len(buf))
	cfb.XORKeyStream(val, buf)

	return append(val, iv...), nil
}

func (p *Settings) decrypt(buf []byte) ([]byte, error) {
	bln := len(buf)
	cln := bln - aes.BlockSize
	ct := buf[0:cln]
	iv := buf[cln:bln]

	cfb := cipher.NewCFBDecrypter(p.cip, iv)
	val := make([]byte, cln)
	cfb.XORKeyStream(val, ct)
	return val, nil
}
