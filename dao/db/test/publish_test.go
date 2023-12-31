package test

import (
	"context"
	"fmt"
	"testing"
	"time"
	"wizh/dao/db"
)

func TestCreateVideo(t *testing.T) {
	video := &db.Video{
		PlayUrl:   "test",
		CoverUrl:  "test",
		Title:     "test video",
		AuthorID:  3,
		CreatedAt: time.Now(),
	}

	if err := db.CreateVideo(context.Background(), video); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("CreateVideo success!")
}

func TestGetVideoListByUserId(t *testing.T) {
	videoList, err := db.GetVideoListByUserId(context.Background(), 3)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("GetVideoListByUserId success!")
	for _, v := range videoList {
		fmt.Println(v)
	}
}

func TestDeleteVideoById(t *testing.T) {
	if err := db.DeleteVideoById(context.Background(), 7, 3); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("DeleteVideoById success!")
}
