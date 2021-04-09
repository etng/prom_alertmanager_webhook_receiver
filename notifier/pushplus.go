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

type PushPlusNotifier struct {
	defaultToken string
	defaultGroup string
}

func NewPushPlusNotifier(defaultToken, defaultGroup string) *PushPlusNotifier {
	return &PushPlusNotifier{defaultToken: defaultToken, defaultGroup: defaultGroup}
}

// TransformToMarkdown transform alertmanager notification to dingtalk markdow message
func (n *PushPlusNotifier) transform(notification model.Notification) (markdown *model.PushPlusMarkdown, robotURL string, err error) {
	annotations := notification.CommonAnnotations
	var topic = annotations["pushplusTopic"]
	var token = annotations["pushplusToken"]
	if token == "" {
		token = n.defaultToken
	}
	if topic == "" {
		topic = n.defaultGroup
	}
	if token == "" {
		return nil, "", fmt.Errorf("no pushplus token")
	}
	robotURL = fmt.Sprintf("http://www.pushplus.plus/send?topic=%s", topic)
	title, body := BuildMarkdown(notification)
	// fmt.Println(body)
	markdown = &model.PushPlusMarkdown{
		Token:    token,
		Template: "markdown",
		Title:    title,
		Content:  body,
	}

	return
}

// Send send markdown message to dingtalk
func (n *PushPlusNotifier) Send(notification model.Notification) (err error) {

	markdown, pushPlusRobotURL, err := n.transform(notification)

	if err != nil {
		return
	}

	data, err := json.Marshal(markdown)
	if err != nil {
		return
	}

	if pushPlusRobotURL == "" {
		return nil
	}
	log.Printf("pushPlusRobotURL: %s", pushPlusRobotURL)

	req, err := http.NewRequest(
		"POST",
		pushPlusRobotURL,
		bytes.NewBuffer(data))

	if err != nil {
		log.Printf("pushPlus robot url notify request build fail: %s", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()
	log.Printf("pp notify response Status: %s", resp.Status)
	log.Printf("pp notify response Headers: %s", resp.Header)
	if resp.StatusCode >= 400 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("pp notify response Body: %s", body)
	}
	return
}
