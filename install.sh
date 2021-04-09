#!/bin/sh -ex
UPDATE_SERVER=http://192.168.4.203:8018
wget -O/etc/systemd/system/prom_am_webhook_receiver.service ${UPDATE_SERVER}/prom_am_webhook_receiver.service
wget -O/usr/local/bin/prom_am_webhook_receiver ${UPDATE_SERVER}/prom_am_webhook_receiver
chmod +x /usr/local/bin/prom_am_webhook_receiver

systemctl enable prom_am_webhook_receiver
systemctl restart prom_am_webhook_receiver
systemctl status prom_am_webhook_receiver
