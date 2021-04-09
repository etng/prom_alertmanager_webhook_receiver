package notifier

import (
	"bytes"
	"fmt"
	"time"

	"github.com/etng/prom_alertmanager_webhook_receiver/model"
)

func strInList(s string, l []string) bool {
	for _, v := range l {
		if s == v {
			return true
		}
	}
	return false
}
func init() {
	messages = map[string]string{
		"alertname": "报警名称",
		"job":       "监控任务",
		"severity":  "严重等级",
		"account":   "帐号",
	}
}

var messages map[string]string

func translate(msgId string) string {
	if msg, ok := messages[msgId]; ok {
		return msg
	}
	return msgId
}
func BuildMarkdown(notification model.Notification) (title string, body string) {
	// groupKey := strings.TrimPrefix(notification.GroupKey, "{}:")
	status := notification.Status
	var buffer bytes.Buffer
	title = fmt.Sprintf("新报警 %s (状态: %s)", notification.GroupLabels["alertname"], status)

	buffer.WriteString(fmt.Sprintf("### %s \n", title))
	blackList := []string{"alertname"}
	buffer.WriteString("\n")
	for k, v := range notification.GroupLabels {
		if strInList(k, blackList) {
			continue
		}
		buffer.WriteString(fmt.Sprintf("* `%s`: `%v`\n", translate(k), v))
		blackList = append(blackList, k)
	}
	for k, v := range notification.CommonLabels {
		if strInList(k, blackList) {
			continue
		}
		buffer.WriteString(fmt.Sprintf("* `%s`: `%v`\n", translate(k), v))
		blackList = append(blackList, k)
	}
	buffer.WriteString("\n")

	buffer.WriteString(fmt.Sprintf("#### 告警项:\n"))
	loc, _ := time.LoadLocation("Asia/Shanghai")
	for i, alert := range notification.Alerts {
		annotations := alert.Annotations
		buffer.WriteString(fmt.Sprintf("##### %d. %s\n > %s\n", i+1, annotations["summary"], annotations["description"]))
		buffer.WriteString(fmt.Sprintf("\n> 开始时间：%s\n", alert.StartsAt.In(loc).Format("15:04:05")))
	}
	body = buffer.String()
	return
}
