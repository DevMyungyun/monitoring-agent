module crontab

go 1.15

require (
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.7.0
	monitoring-agent/command v0.0.0
	monitoring-agent/http v0.0.0
)

replace (
	monitoring-agent/command v0.0.0 => ./command
	monitoring-agent/http v0.0.0 => ./http
)
