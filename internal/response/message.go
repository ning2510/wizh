package response

import "wizh/kitex/kitex_gen/message"

type MessageAction struct {
	Base
}

type MessageChat struct {
	Base
	MessageList []*message.Message `json:"message_list"`
}
