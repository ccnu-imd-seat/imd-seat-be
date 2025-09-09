package image

import (
	"bytes"
	"context"
	"imd-seat-be/internal/config"

	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
)

type Uploader interface {
	Upload(ctx context.Context, data []byte, key, FileName string) (string, error)
	GetImageURL(key string) string
}

type QiniuUploader struct {
	cfg config.Config
	uploaderManager *uploader.UploadManager
}

func NewQiniuUploader(cfg config.Config) *QiniuUploader {
	mac := credentials.NewCredentials(cfg.Qiniu.AccessKey, cfg.Qiniu.SecretKey)

	uploaderManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})
	return &QiniuUploader{
		cfg: cfg,
		uploaderManager: uploaderManager,
	}
}

func (q *QiniuUploader) Upload(ctx context.Context, data []byte, key, FileName string) (string, error) {
	err := q.uploaderManager.UploadReader(ctx, bytes.NewReader(data), &uploader.ObjectOptions{
		BucketName: q.cfg.Qiniu.Bucket,
		ObjectName: &key,
		FileName:   FileName,
		CustomVars: map[string]string{
			"name": "qrcode",
		},
	}, nil)
	if err != nil {
		return "", err
	}
	return q.GetImageURL(key), nil
}

// 获取url
func (q *QiniuUploader) GetImageURL(key string) string {
	return q.cfg.Qiniu.Domain + "/" + key
}
