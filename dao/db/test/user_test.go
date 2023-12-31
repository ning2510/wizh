package test

import (
	"context"
	"fmt"
	"testing"
	"wizh/dao/db"
)

func TestGetUsersByIds(t *testing.T) {
	userIds := []int64{1, 2, 3}
	users, err := db.GetUsersByIds(context.Background(), userIds)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("GetUsersByIds success!")
	for _, u := range users {
		fmt.Println(u)
	}
}

func TestGetUserById(t *testing.T) {
	user, err := db.GetUserById(context.Background(), 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("GetUserById success!")
	fmt.Println(user)
}

func TestGetUserByUsername(t *testing.T) {
	user, err := db.GetUserByUsername(context.Background(), "ccc")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("GetUserByUsername success!")
	fmt.Println(user)
}

func TestCreateUsers(t *testing.T) {
	userList := make([]*db.User, 0)
	for i := 0; i < 2; i++ {
		userList = append(userList, &db.User{
			UserName: fmt.Sprintf("TestUser%d", i),
			Password: "123",
		})
	}

	if err := db.CreateUsers(context.Background(), userList); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("CreateUsers success!")
}

func TestCreateUser(t *testing.T) {
	user := &db.User{
		UserName: "ccc",
		Password: "123",
	}

	if err := db.CreateUser(context.Background(), user); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("CreateUser success!")
}

func TestDeleteUserById(t *testing.T) {
	if err := db.DeleteUserById(context.Background(), 1); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("DeleteUserById success!")
}

func TestDeleteUserByIds(t *testing.T) {
	if err := db.DeleteUserByIds(context.Background(), []int64{4, 5}); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("DeleteUserByIds success!")
}
