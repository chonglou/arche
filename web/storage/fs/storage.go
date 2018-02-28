package fs

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/chonglou/arche/web/storage"
	"github.com/google/uuid"
)

// New new s3
func New(root, endpoint string) (storage.Storage, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, err
	}
	return &Fs{root: root, endpoint: endpoint}, nil
}

// Fs file-system
type Fs struct {
	root     string
	endpoint string
}

// Write write to
func (p *Fs) Write(name string, body []byte, size int64) (string, string, error) {
	fn := uuid.New().String() + filepath.Ext(name)
	dest := path.Join(p.root, fn)
	fileType := http.DetectContentType(body)
	if err := ioutil.WriteFile(dest, body, 0644); err != nil {
		return "", "", err
	}
	return fileType, p.endpoint + "/" + fn, nil
}
