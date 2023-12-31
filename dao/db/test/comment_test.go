package test

import (
	"context"
	"fmt"
	"testing"
	"time"
	"wizh/dao/db"
)

func TestCreateComment(t *testing.T) {
	comment := &db.Comment{
		Content:   "test comment",
		VideoID:   2,
		UserID:    3,
		CreatedAt: time.Now(),
	}

	if err := db.CreateComment(context.Background(), comment); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("CreateComment success!")
}

func TestDeleteCommentById(t *testing.T) {
	if err := db.DeleteCommentById(context.Background(), 1, 2); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("DeleteCommentById success!")
}

func TestGetCommentListByVideoId(t *testing.T) {
	commentList, err := db.GetCommentListByVideoId(context.Background(), 2)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("GetCommentListByVideoId success!")
	for _, c := range commentList {
		fmt.Println(c)
	}
}

func TestGetCommentByCommentId(t *testing.T) {
	comment, err := db.GetCommentByCommentId(context.Background(), 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("GetCommentByCommentId success!")
	fmt.Println(comment)
}
