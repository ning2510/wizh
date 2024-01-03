package test

import (
	"context"
	"fmt"
	"testing"
	"wizh/dao/db"
)

func TestCreateVideoFavorite(t *testing.T) {
	if err := db.CreateVideoFavorite(context.Background(), 3, 2, 3); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("CreateVideoFavorite success!")
}

func TestDeleteFavoriteVideoByUserVideoId(t *testing.T) {
	if err := db.DeleteFavoriteVideoByUserVideoId(context.Background(), 3, 2, 3); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("DeleteFavoriteVideoByUserVideoId success!")
}

func TestGetFavoriteVideoRelationByUserVideoId(t *testing.T) {
	favoriteVideoRelation, err := db.GetFavoriteVideoRelationByUserVideoId(context.Background(), 3, 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("GetFavoriteVideoRelationByUserVideoId success!")
	fmt.Println(favoriteVideoRelation.UserID, favoriteVideoRelation.VideoID)
}

func TestGetFavoriteListByUserId(t *testing.T) {
	favoriteList, err := db.GetFavoriteListByUserId(context.Background(), 3)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("GetFavoriteListByUserId success!")
	for _, v := range favoriteList {
		fmt.Println(v.UserID, v.VideoID)
	}
}

func TestGetAllFavoriteList(t *testing.T) {
	favoriteList, err := db.GetAllFavoriteList(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("GetAllFavoriteList success!")
	for _, v := range favoriteList {
		fmt.Println(v.UserID, v.VideoID)
	}
}
