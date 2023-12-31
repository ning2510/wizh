package test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"testing"
	"time"
	"wizh/pkg/minio"

	"github.com/gin-gonic/gin"
)

func ExpectEqual(left interface{}, right interface{}, t *testing.T) {
	if left != right {
		t.Fatalf("expected %v == %v\n%s", left, right, debug.Stack())
	}
}

func TestUploadFileByPath(t *testing.T) {
	bucketName := "test-minio"
	if err := minio.CreateBucket(bucketName); err != nil {
		panic(err)
	}

	objectName := "minio-test.txt"
	filePath := "test.txt"
	contentType := "application/txt"

	// 检查文件是否存在并获取其大小
	file, err := os.Open(filePath)
	ExpectEqual(err, nil, t)
	defer file.Close()

	fileStat, err := file.Stat()
	ExpectEqual(err, nil, t)

	size, err := minio.UploadFileByPath(bucketName, objectName, filePath, contentType)

	ExpectEqual(size, fileStat.Size(), t)
	ExpectEqual(err, nil, t)
}

func TestUploadFileByIO(t *testing.T) {
	bucketName := "test-minio"
	if err := minio.CreateBucket(bucketName); err != nil {
		panic(err)
	}

	objectName := "minio-test2.txt"
	filePath := "test.txt"
	contentType := "application/txt"

	// 搭建一个简单的 web 服务器用于接收文件
	r := gin.Default()
	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		ExpectEqual(err, nil, t)
		fp, err := file.Open()
		ExpectEqual(err, nil, t)
		size, err := minio.UploadFileByIO(bucketName, objectName, fp, file.Size, contentType)
		ExpectEqual(size, file.Size, t)
		ExpectEqual(err, nil, t)
	})

	go func() {
		r.Run("0.0.0.0:9999")
	}()

	// 使用 http client 上传文件
	time.Sleep(5 * time.Second)
	file, err := os.Open(filePath)
	ExpectEqual(err, nil, t)
	defer file.Close()

	content := &bytes.Buffer{}
	writer := multipart.NewWriter(content)
	form, err := writer.CreateFormFile("file", filepath.Base(filePath))
	ExpectEqual(err, nil, t)

	_, err = io.Copy(form, file)
	ExpectEqual(err, nil, t)

	err = writer.Close()
	ExpectEqual(err, nil, t)

	req, err := http.NewRequest("POST", "http://10.13.0.38:9999/upload", content)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	ExpectEqual(err, nil, t)

	client := &http.Client{}
	_, err = client.Do(req)
	ExpectEqual(err, nil, t)
}

func TestGetFileTemporaryURL(t *testing.T) {
	bucketName := "test-minio"
	objectName := "minio-test.txt"

	url, err := minio.GetFileTemporaryURL(bucketName, objectName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(url)
}

func TestRemoveBucket(t *testing.T) {
	if err := minio.RemoveBucket(context.Background(), "test-minio"); err != nil {
		fmt.Printf("RemoveBucket error! %v\n", err)
		return
	}
	fmt.Println("RemoveBucket success!")
}
