SHELL := /bin/bash

DEST=root@192.168.4.203:/home/static/
DEMO_HOST=localhost:9095
# DEMO_HOST=192.168.4.203:9095
run:
	source .env && echo go run cmd/webhook/webhook.go -token=$${DingTalkToken} -prefix=$${DingTalkPrefix} -port=$${ListenPort}  -pp_token=$${PushPlusToken} -pp_topic=$${PushPlusTopic} -log=webhook.log -debug
	source .env && go run cmd/webhook/webhook.go -token=$${DingTalkToken} -prefix=$${DingTalkPrefix} -port=$${ListenPort}  -pp_token=$${PushPlusToken} -pp_topic=$${PushPlusTopic} -log=webhook.log -debug
mock:
	curl -XPOST -H "Content-Type: application/json" ${DEMO_HOST}/webhook -d @sample_alert_request.json
build:
	go build -o prom_am_webhook_receiver cmd/webhook/webhook.go
	upx prom_am_webhook_receiver
deploy: build
	scp prom_am_webhook_receiver ${DEST}
	scp prom_am_webhook_receiver.service ${DEST}
	scp install.sh ${DEST}/prom_am_webhook_receiver.install.sh