package model

type DingTalkMessage struct {
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type DingTalkMarkdown struct {
	MsgType  string    `json:"msgtype"`
	At       *At       `json:"at"`
	Markdown *Markdown `json:"markdown"`
}

type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type PushPlusMarkdown struct {
	Template string `json:"template"`
	Title    string `json:"title"`

	Token   string `json:"token"`
	Content string `json:"content"`
}
