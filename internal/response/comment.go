package response

import "wizh/kitex/kitex_gen/comment"

type CommentAction struct {
	Base
}

type CommentList struct {
	Base
	CommentList []*comment.Comment `json:"comment_list"`
}
