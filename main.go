package main

import (
	// "reflect"
	"fmt"
	"io/ioutil"
	"net/http"
	// "bytes"
	"encoding/json"
	"crypto/aes"
	"strings"
	"time"

	"monitoring-agent/command"
	encryption "monitoring-agent/encryption"
	jwt "monitoring-agent/auth"
	// httpReqRes "monitoring-agent/http"
	"monitoring-agent/crontab"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	
	var checkOS = command.DetectOS()
	var resource = command.GetResource(checkOS)	
	crontab.Intialize("@every 0h0m30s", checkOS, resource)
}

func checkError(err error) {
	if err != nil {
	  fmt.Println(err)
	  os.Exit(1)
	}
}

func main() {
	r := gin.Default()

	r.POST("/handshake", handshake)

	v1 := r.Group("/v1")
	{
		v1.GET("/health", health)
		v1.POST("/cron/start", cronStart)
		v1.GET("/cron/check", checkCron)
		v1.GET("/cron/stop", cronStop)
		v1.GET("/token/check", tokenCheck)
	}
	r.Run(":9000")
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"code": 200,
	})
}

func handshake(c *gin.Context) {
	reqBody, _ := ioutil.ReadAll(c.Request.Body)
	
	type agentData struct{
		ID string
		NAME string
		JWT string
		MAIN_SERVER_ADDRESS string
		OS string
	}
	var ad agentData
	trimStrBody := strings.Replace(string(reqBody), " ", "", -1)
	trimStrBody = strings.Replace(trimStrBody, "\n", "", -1)
	err := json.Unmarshal([]byte(trimStrBody), &ad)
	if err != nil {
		defer func() { 
			s := recover() 
			fmt.Println(s)
		}()
        panic(err)
	}

	checkOS := command.DetectOS()
	path := ""
	switch checkOS {
	case "windows":
		path = "C:\\temp\\agent-config"
	case "darwin":
		path = "/tmp/agent-config"
	case "linux":
		path = "/tmp/agent-config"
	default:
		fmt.Println("This OS is not supported : ", checkOS)
	}
	config := make(map[string]string)
	config["ID"] = ad.ID
	config["NAME"] = ad.NAME
	config["JWT"] = ad.JWT
	config["MAIN_SERVER_ADDRESS"] = ad.MAIN_SERVER_ADDRESS
	config["CREATE_AT"] = time.Now().String()

	bodyBytes, _ := json.Marshal(config)

	key := "16byteSecret!!!!" // must be 16 byte
	block, err := aes.NewCipher([]byte(key))
	if err != nil {	
		fmt.Println(err)
		return
	}

	ciphertext := encryption.Encrypt(block, []byte(bodyBytes))
	fmt.Printf("%x\n", ciphertext)

    err = ioutil.WriteFile(path, ciphertext, 0644)
    if err != nil {
        panic(err)
    }
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func cronStart(c *gin.Context) {
	payload := jwt.TokenValidCheck(c.Request)
	log.Info("payload####",payload)
	log.Info("Create new cron")
	// Start cron with one scheduled job
	log.Info("Start cron")
	crontab.Start()
	// time.Sleep(2 * time.Minute)
	log.Infof("cron entity : %+v\n", crontab.Entries())

	c.JSON(http.StatusOK, gin.H{"message": "This agent start to work..."})
}

func checkCron(c *gin.Context) {
	payload := jwt.TokenValidCheck(c.Request)
	log.Info("payload####",payload)
	log.Info("Check cron")
	// println(Cron.Entry(1))
	c.JSON(http.StatusOK, gin.H{"message": "Check this agent work..."})
}

func cronStop(c *gin.Context) {
	payload := jwt.TokenValidCheck(c.Request)
	log.Info("payload####",payload)
	log.Info("Stop cron")
	c.JSON(http.StatusOK, gin.H{"message": "Stop this agent work..."})
	log.Infof("cron entity : %+v\n", crontab.Entries())
	crontab.Stop()	
	// Cron.Remove(1)
}

func tokenCheck(c *gin.Context) {
	payload := jwt.TokenValidCheck(c.Request)
	log.Info("payload####",payload)
	c.JSON(http.StatusOK, gin.H{"message": "jwt token check..."})
}