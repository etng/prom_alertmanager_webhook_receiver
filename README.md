## Prometheus Alertmanager Webhook Receiver

Webhook service support send Prometheus 2.0 alert message to Robot or Push services.

## Usage

```shell
go build -o prom_am_webhook_receiver cmd/webhook/webhook.go
./prom_am_webhook_receiver -token=xxxx -prefix= -port=9095
```

* -token: default dingtalk robot token
* -prefix: default dingtalk robot prefix

Or you can overwrite by add annotations to Prometheus alertrule to special the webhook for each alert rule.

```yaml
groups:
- name: hostStatsAlert
  rules:
  - alert: hostCpuUsageAlert
    expr: sum(avg without (cpu)(irate(node_cpu_seconds_total{mode!='idle'}[5m]))) by (instance) > 0.85
    for: 1m
    labels:
      severity: page
    annotations:
      summary: "Instance {{ $labels.instance }} CPU usgae high"
      description: "{{ $labels.instance }} CPU usage above 85% (current value: {{ $value }})"
      dingtalkToken: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
      dingtalkPrefix: "打扰一下："
```
