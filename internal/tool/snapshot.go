package tool

import (
	"bytes"
	"fmt"
	"os"
	"wizh/pkg/zap"

	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func GetSnapshot(videoPath, snapshotPath string, frameNum int) (string, error) {
	logger := zap.InitLogger()
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(videoPath).Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()

	if err != nil {
		logger.Errorf("生成缩略图失败: %v\n", err)
		return "", err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		logger.Errorf("生成缩略图失败: %v\n", err)
		return "", err
	}

	err = imaging.Save(img, snapshotPath+".png")
	if err != nil {
		logger.Errorf("生成缩略图失败: %v\n", err)
		return "", err
	}

	imgPath := snapshotPath + ".png"
	return imgPath, nil
}

func GetSnapshotImageBuffer(videoPath string, frameNum int) (*bytes.Buffer, error) {
	logger := zap.InitLogger()
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(videoPath).Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()

	if err != nil {
		logger.Errorf("生成缩略图失败: %v\n", err.Error())
		return nil, err
	}
	return buf, nil
}
