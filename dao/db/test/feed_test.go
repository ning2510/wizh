package test

import (
	"context"
	"fmt"
	"testing"
	"time"
	"wizh/dao/db"
)

func TestGetVideos(t *testing.T) {
	latestTime := time.Now().UnixMilli()
	videoList, err := db.GetVideos(context.Background(), 10, &latestTime)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("GetVideos success!")
	for _, v := range videoList {
		fmt.Println(v)
	}
}

func TestGetVideoByVideoId(t *testing.T) {
	video, err := db.GetVideoByVideoId(context.Background(), 5)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("GetVideoByVideoId success!")
	fmt.Println(video)
}

func TestGetVideosByVideoIds(t *testing.T) {
	videoIds := []int64{2}
	videoList, err := db.GetVideosByVideoIds(context.Background(), videoIds)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("GetVideosByVideoIds success!")
	for _, v := range videoList {
		fmt.Println(v)
	}
}
