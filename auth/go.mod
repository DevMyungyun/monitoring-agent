module auth

go 1.15

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/sirupsen/logrus v1.7.0
	monitoring-agent/encryption v0.0.0
)

replace (
	monitoring-agent/encryption v0.0.0 => ./encryption
)