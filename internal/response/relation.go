package response

import "wizh/kitex/kitex_gen/user"

type RelationAction struct {
	Base
}

type RelationFollowList struct {
	Base
	UserList []*user.User `json:"user_list"`
}

type RelationFollowerList struct {
	Base
	UserList []*user.User `json:"user_list"`
}
