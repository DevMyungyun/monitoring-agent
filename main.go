package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var Cron = cron.New()

type cronEntities struct {
	id       int
	schedule string
}

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
}

func main() {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.GET("/health", health)
		v1.POST("/signup", signup)
		v1.POST("/login", login)
		v1.POST("/cron/start", cronStart)
		v1.GET("/cron/check", checkCron)
		v1.GET("/cron/stop", cronStop)
	}
	r.Run(":9000")
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func signup(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "signed up",
	})
}

func login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "logged in",
	})
}

func cronStart(c *gin.Context) {
	log.Info("Create new cron")

	Cron.AddFunc("*/1 * * * *", func() { log.Info("[Job 1]Every minute job\n") })

	// Start cron with one scheduled job
	log.Info("Start cron")
	Cron.Start()
	printCronEntries(Cron.Entries())
	// time.Sleep(2 * time.Minute)

	c.JSON(http.StatusOK, gin.H{"message": "This agent start to work..."})

}

func checkCron(c *gin.Context) {
	log.Info("Check cron")

	c.JSON(http.StatusOK, gin.H{"message": "Check this agent work..."})
	println(Cron.Entries())
	// println(Cron.Parse())
	// var Cron = cron.New()
	// inspect(Cron.Entries())
}

func cronStop(c *gin.Context) {
	log.Info("Stop cron")

	c.JSON(http.StatusOK, gin.H{"message": "Stop this agent work..."})

	// var Cron = cron.New()
	Cron.Stop()
}

func printCronEntries(cronEntries []cron.Entry) {
	log.Infof("Cron Info: %+v\n", cronEntries)
}
