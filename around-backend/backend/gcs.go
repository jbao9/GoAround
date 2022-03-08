package backend

//GCS library
//https://github.com/GoogleCloudPlatform/golang-samples/blob/master/storage/objects/main.go

import (
	"context"
	"fmt"
	"io"

	"around/util"

	"cloud.google.com/go/storage"
)

var (
	GCSBackend *GoogleCloudStorageBackend
)

type GoogleCloudStorageBackend struct {
	client *storage.Client
	bucket string
}

func InitGCSBackend(config *util.GCSInfo) {
	client, err := storage.NewClient(context.Background())
	//client: sessionFactory
	if err != nil {
		panic(err)
	}

	GCSBackend = &GoogleCloudStorageBackend{
		client: client,
		bucket: config.Bucket,
	}
}

func (backend *GoogleCloudStorageBackend) SaveToGCS(r io.Reader, objectName string) (string, error) {
	//backend 相当于java的 this														//返回类型1，用户上传文件的rul 返回类型2
	ctx := context.Background()
	object := backend.client.Bucket(backend.bucket).Object(objectName)
	wc := object.NewWriter(ctx)
	if _, err := io.Copy(wc, r); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	//给前端打开读取文件的权限
	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		//ACL: access control list    //set all users have the access of read 上传到网盘里的文件读权限要向所有人敞开
		return "", err
	}

	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}

	fmt.Printf("File is saved to GCS: %s\n", attrs.MediaLink)
	return attrs.MediaLink, nil

}
