package s3

import (
	"bytes"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/chonglou/arche/web/storage"
	"github.com/google/uuid"
)

// New new s3
func New(accessKeyID, secretAccessKey, region, bucket string) (storage.Storage, error) {
	creds := credentials.NewStaticCredentials(accessKeyID, secretAccessKey, "")
	if _, err := creds.Get(); err != nil {
		return nil, err
	}
	return &S3{credentials: creds, region: region, bucket: bucket}, nil
}

// S3 amazon s3
type S3 struct {
	region      string
	bucket      string
	credentials *credentials.Credentials
}

// Write write to
func (p *S3) Write(name string, body []byte, size int64) (string, string, error) {

	svc := s3.New(
		session.New(),
		aws.NewConfig().WithRegion(p.region).WithCredentials(p.credentials),
	)

	fn := "/upload/" + uuid.New().String() + filepath.Ext(name)

	fileBytes := bytes.NewReader(body)
	fileType := http.DetectContentType(body)

	params := &s3.PutObjectInput{
		ACL:           aws.String("public-read"),
		Bucket:        aws.String(p.bucket),
		Key:           aws.String(fn),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}

	_, err := svc.PutObject(params)
	if err != nil {
		return "", "", err
	}

	href := "https://s3-" + p.region + ".amazonaws.com/" + p.bucket + fn // FIXME
	return fileType, href, nil

}
