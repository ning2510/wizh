package response

import "wizh/kitex/kitex_gen/video"

type FavoriteVideoAction struct {
	Base
}

type FavoriteVideoList struct {
	Base
	VideoList []*video.Video `json:"video_list"`
}

type FavoriteCommentAction struct {
	Base
}
