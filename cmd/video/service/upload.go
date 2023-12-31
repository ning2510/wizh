package service

import (
	"bytes"
	"wizh/internal/tool"
	"wizh/pkg/minio"
	"wizh/pkg/zap"
)

// 上传视频获取封面
func VideoPublish(data []byte, videoTitle string, coverTitle string) error {
	playUrl, err := uploadVideo(data, videoTitle)
	if err != nil {
		return err
	}

	_, err = uploadCover(playUrl, coverTitle)
	if err != nil {
		return err
	}
	return nil
}

// 上传视频 minio
func uploadVideo(data []byte, videoTitle string) (string, error) {
	logger := zap.InitLogger()
	reader := bytes.NewReader(data)
	contentType := "application/mp4"

	uploadSize, err := minio.UploadFileByIO(minio.VideoBucketName, videoTitle, reader, int64(len(data)), contentType)
	if err != nil {
		logger.Errorf("上传视频到 minio 失败: %v\n", err)
		return "", err
	}
	logger.Infof("视频文件大小为: %v", uploadSize)

	playUrl, err := minio.GetFileTemporaryURL(minio.VideoBucketName, videoTitle)
	if err != nil {
		logger.Errorf("获取视频链接失败: %v\n", err)
		return "", err
	}
	return playUrl, nil
}

// 截取封面并上传到 minio
func uploadCover(playUrl string, coverTitle string) (string, error) {
	logger := zap.InitLogger()
	imgBuffer, err := tool.GetSnapshotImageBuffer(playUrl, 1)
	if err != nil {
		logger.Errorf("服务器内部错误，封面获取失败: %v\n", err)
		return "", err
	}

	var imgByte []byte
	imgBuffer.Write(imgByte)
	contentType := "application/png"

	uploadSize, err := minio.UploadFileByIO(minio.CoverBucketName, coverTitle, imgBuffer, int64(imgBuffer.Len()), contentType)
	if err != nil {
		logger.Errorf("上传封面到 minio 失败: %v\n", err)
		return "", err
	}
	logger.Infof("封面文件大小为: %v", uploadSize)

	coverUrl, err := minio.GetFileTemporaryURL(minio.CoverBucketName, coverTitle)
	if err != nil {
		logger.Errorf("获取封面链接失败: %v\n", err)
		return "", err
	}
	return coverUrl, nil
}
