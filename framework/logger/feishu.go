package logger

//飞书消息结构
type FeishuTitleContent struct {
	//Content gitlab.MergeRequestEventPayload `json:"content"`
	Tag  string `json:"tag"`
	Text string `json:"text"`
}

type FeishuTitle struct {
	//Content gitlab.MergeRequestEventPayload `json:"content"`
	Title   string                 `json:"title"`
	Content [][]FeishuTitleContent `json:"content"`
}

type FeishuZh_cn struct {
	Zh_cn FeishuTitle `json:"zh_cn"`
}

type FeishuContent struct {
	Post FeishuZh_cn `json:"post"`
}

type FeishuMsg struct {
	MsgType string        `json:"msg_type"`
	Content FeishuContent `json:"content"`
}
