module monitoring-agent

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.7.0
	github.com/thoas/go-funk v0.7.0
	monitoring-agent/auth v0.0.0
	monitoring-agent/command v0.0.0
	monitoring-agent/encryption v0.0.0
	monitoring-agent/http v0.0.0

)

replace (
	monitoring-agent/auth v0.0.0 => ./auth
	monitoring-agent/command v0.0.0 => ./command
	monitoring-agent/encryption v0.0.0 => ./encryption
	monitoring-agent/http v0.0.0 => ./http
)
