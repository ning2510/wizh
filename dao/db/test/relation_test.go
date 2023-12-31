package test

import (
	"context"
	"fmt"
	"testing"
	"wizh/dao/db"
)

func TestGetRelationByUserIds(t *testing.T) {
	followRelation, err := db.GetRelationByUserIds(context.Background(), 3, 6)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("GetRelationByUserIds success!")
	fmt.Println(followRelation.UserID)
	fmt.Println(followRelation.User)
	fmt.Println(followRelation.ToUserID)
	fmt.Println(followRelation.ToUser)
}

func TestCreateRelation(t *testing.T) {
	if err := db.CreateRelation(context.Background(), 3, 6); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("CreateRelation success!")
}

func TestDeleteRelationByUserIds(t *testing.T) {
	err := db.DeleteRelationByUserIds(context.Background(), 3, 6)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("DeleteRelationByUserIds success!")
}

func TestGetFollowingListByUserId(t *testing.T) {
	followRelationList, err := db.GetFollowingListByUserId(context.Background(), 3)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("GetFollowingListByUserId success!")
	for _, f := range followRelationList {
		fmt.Println(f)
	}
}

func TestGetFollowerListByUserId(t *testing.T) {
	followerRelationList, err := db.GetFollowerListByUserId(context.Background(), 6)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("GetFollowerListByUserId success!")
	for _, f := range followerRelationList {
		fmt.Println(f)
	}
}
