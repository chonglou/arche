package nut

import (
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"sync"

	"github.com/astaxie/beego"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/google/uuid"
)

var (
	_storage     Storage
	_storageOnce sync.Once
)

// STORAGE storage instance
func STORAGE() Storage {
	_storageOnce.Do(func() {
		prv := beego.AppConfig.String("storageprovider")
		switch prv {
		case "s3":
			creds := credentials.NewStaticCredentials(
				beego.AppConfig.String("awsaccesskeyid"),
				beego.AppConfig.String("awssecretaccesskey"),
				"",
			)
			if _, err := creds.Get(); err != nil {
				beego.Error(err)
				return
			}
			_storage = &S3Storage{
				credentials: creds,
				region:      beego.AppConfig.String("awss3region"),
				bucket:      beego.AppConfig.String("awss3bucket"),
			}
		case "local":
			_storage = &FileSystemStorage{
				root: beego.AppConfig.String("localstoragedir"),
				home: beego.AppConfig.String("localstorageendpoint"),
			}
		default:
			beego.Error("storage provider", prv, "is not support")
		}
	})
	return _storage
}

// Storage storage
type Storage interface {
	Save(name string, body []byte, size int64) (fileType string, url string, err error)
}

// FileSystemStorage file system storage
type FileSystemStorage struct {
	root string
	home string
}

// Save write to
func (p *FileSystemStorage) Save(name string, body []byte, size int64) (string, string, error) {
	fn := uuid.New().String() + filepath.Ext(name)
	dest := path.Join(p.root, fn)
	beego.Debug("generate file", dest)
	fileType := http.DetectContentType(body)
	if err := ioutil.WriteFile(dest, body, 0644); err != nil {
		return "", "", err
	}
	return fileType, p.home + "/" + fn, nil
}
