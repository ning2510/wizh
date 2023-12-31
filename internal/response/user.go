package response

import "wizh/kitex/kitex_gen/user"

type Register struct {
	Base
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type Login struct {
	Base
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserInfo struct {
	Base
	User *user.User `json:"user"`
}
