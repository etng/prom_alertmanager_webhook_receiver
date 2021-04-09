package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/etng/prom_alertmanager_webhook_receiver/model"
)

type Notifier interface {
	Send(notification model.Notification, options NotifyOptions)
}
type DingTalkNotifier struct {
	defaultToken  string
	defaultPrefix string
}

func NewDingTalkNotifier(defaultToken, defaultPrefix string) *DingTalkNotifier {
	return &DingTalkNotifier{defaultToken: defaultToken, defaultPrefix: defaultPrefix}
}

// TransformToMarkdown transform alertmanager notification to dingtalk markdow message
func (n *DingTalkNotifier) transform(notification model.Notification) (markdown *model.DingTalkMarkdown, robotURL string, err error) {

	annotations := notification.CommonAnnotations
	var dingtalkToken = annotations["dingtalkToken"]
	var dingtalkPrefix = annotations["dingtalkPrefix"]
	if dingtalkToken == "" {
		dingtalkToken = n.defaultToken
	}
	if dingtalkPrefix == "" {
		dingtalkPrefix = n.defaultPrefix
	}
	if dingtalkToken == "" {
		return nil, "", fmt.Errorf("no dingtalk token")
	}
	robotURL = fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", dingtalkToken)

	title, body := BuildMarkdown(notification)
	if dingtalkPrefix != "" {
		body = fmt.Sprintf("%s\n%s", dingtalkPrefix, body)
	}

	markdown = &model.DingTalkMarkdown{
		MsgType: "markdown",
		Markdown: &model.Markdown{
			Title: title,
			Text:  body,
		},
		At: &model.At{
			IsAtAll: false,
		},
	}

	return
}

// Send send markdown message to dingtalk
func (n *DingTalkNotifier) Send(notification model.Notification) (err error) {

	markdown, dingtalkRobotURL, err := n.transform(notification)

	if err != nil {
		return
	}

	data, err := json.Marshal(markdown)
	if err != nil {
		return
	}

	if dingtalkRobotURL == "" {
		return nil
	}
	log.Printf("dingtalkRobotURL: %s", dingtalkRobotURL)

	req, err := http.NewRequest(
		"POST",
		dingtalkRobotURL,
		bytes.NewBuffer(data))

	if err != nil {
		log.Printf("dingtalk robot url notify request build fail: %s", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()
	log.Printf("dingding notify response Status: %s", resp.Status)
	log.Printf("dingding notify response Headers: %s", resp.Header)
	log.Printf("dingding Bxpunish flag is %q", resp.Header.Get("Bxpunish"))
	if resp.StatusCode >= 400 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("dingding notify response Body: %s", body)
	}
	return
}
