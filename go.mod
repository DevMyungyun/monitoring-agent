module monitoring-agent

go 1.14

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.6.0
	github.com/thoas/go-funk v0.7.0
	monitoring-agent/command v0.0.0

)

replace (
	monitoring-agent/command v0.0.0 => ./command
)
