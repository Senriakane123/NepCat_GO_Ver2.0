package MSGModel

type ResMessage struct {
	PostType      string  `json:"post_type"`
	SelfID        int64   `json:"self_id"`
	UserID        int64   `json:"user_id,omitempty"`
	RawMessage    string  `json:"raw_message,omitempty"`
	GroupID       int64   `json:"group_id,omitempty"`
	Time          int64   `json:"time"`
	MetaEventType string  `json:"meta_event_type,omitempty"`
	Sender        SendDer `json:"sender"`
	MessageID     int64   `json:"message_id"`
	MessageSeq    int64   `json:"message_seq"`
	MessageType   string  `json:"message_type"`       // "private" 或 "group"
	SubType       string  `json:"sub_type,omitempty"` // 私聊才有 "friend"
}

// 消息发送者信息结构
type SendDer struct {
	UserID   int64  `json:"user_id"`
	NickName string `json:"nickname"`
	Card     string `json:"card"`
	Role     string `json:"role"`
}
