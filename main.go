package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var Cron = cron.New()

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	Cron.AddFunc("*/1 * * * *", func() { log.Info("[Job 1]Every minute job\n") })
}

func main() {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.GET("/health", health)
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

func cronStart(c *gin.Context) {
	log.Info("Create new cron")

	// Start cron with one scheduled job
	log.Info("Start cron")
	Cron.Start()
	printCronEntries(Cron.Entries())
	// time.Sleep(2 * time.Minute)
	println("cron entity :", Cron.Entries())

	c.JSON(http.StatusOK, gin.H{"message": "This agent start to work..."})
}

func checkCron(c *gin.Context) {
	log.Info("Check cron")
	// println(Cron.Entry(1))
	c.JSON(http.StatusOK, gin.H{"message": "Check this agent work..."})
}

func cronStop(c *gin.Context) {
	log.Info("Stop cron")

	c.JSON(http.StatusOK, gin.H{"message": "Stop this agent work..."})

	Cron.Stop()
	// Cron.Remove(1)
}

func printCronEntries(cronEntries []cron.Entry) {
	log.Infof("Cron Info: %+v\n", cronEntries)
}
