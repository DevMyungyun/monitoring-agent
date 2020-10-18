package main

import (
	// "reflect"
	"fmt"
	"io/ioutil"
	"net/http"
	"bytes"
	"encoding/json"
	"crypto/aes"
	"strings"
	"os"
	"time"

	"monitoring-agent/command"
	encryption "monitoring-agent/encryption"
	jwt "monitoring-agent/auth"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var Cron = cron.New()

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	os := command.DetectOS()
	
	cron.WithSeconds()
	// Cron.AddFunc("*/1 * * * *", func() {
	Cron.AddFunc("@every 0h0m30s", func() {
		log.Info("[Job 1]Every minute job\n")
		resource := command.GetResource(os)
		// log.Info(">>>>", reflect.TypeOf(resource))
		// log.Info(">>>>", resource)
		
		res := httpReq("POST", "http://localhost:3333/agent/v1/resource/receive", nil, resource)
		log.Info("response : ", res)
	})
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
	})
}

func handshake(c *gin.Context) {
	reqBody, _ := ioutil.ReadAll(c.Request.Body)
	
	type agentData struct{
		ID string
		JWT string
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
	config["JWT"] = ad.JWT
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
	Cron.Start()
	printCronEntries(Cron.Entries())
	// time.Sleep(2 * time.Minute)
	println("cron entity :", Cron.Entries())

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
	log.Info("Stop cron")

	c.JSON(http.StatusOK, gin.H{"message": "Stop this agent work..."})

	Cron.Stop()
	// Cron.Remove(1)
}

func tokenCheck(c *gin.Context) {
	payload := jwt.TokenValidCheck(c.Request)
	log.Info("payload####",payload)
	c.JSON(http.StatusOK, gin.H{"message": "jwt token check..."})
}

func printCronEntries(cronEntries []cron.Entry) {
	log.Infof("Cron Info: %+v\n", cronEntries)
}

func httpReq(method string, url string, header []string, body interface{}) string {
	bodyBytes, _ := json.Marshal(body)
	bodyBuffer := bytes.NewBuffer(bodyBytes)
	// log.Info("### ", bodyBytes)
	// log.Info("### ", bodyBuffer)
	// Request 객체 생성
	
	req, err := http.NewRequest("POST",url, bodyBuffer)
	if err !=nil {
		panic(err)
	}

	q := req.URL.Query()
    q.Add("name", "777")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", "bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiI3NzciLCJpZCI6ImM2ZDEyMzliLTk4ZTItNDYxYy1hOWQwLTNhNzRmYzhjMzk0OSIsImlhdCI6MTYwMjQ4OTQ4NywiZXhwIjoxNjM0MDI1NDg3fQ.97t7_MWb5ebW_xjIXlDNz28GPQzgxFyG6g7MR2mY5KQhD2C1HHDW6z0BPSkwXQjGqjNy-qb4V6cF_KdSl9WJmA")
	req.Header.Add("Content-Type", "application/json")
	

	// req, err := http.NewRequest(method, url, bodyBuffer)
	// if err != nil {
	// 	panic(err)
	// }
	//필요시 헤더 추가 가능
	// if header != nil {
	// 	for i, v := range header {
	// 		req.Header.Add(string(v[i/2]), string(v[i/2+1]))
	// 	}
	// }
	// Client객체에서 Request 실행
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 결과 출력
	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes) //바이트를 문자열로
	fmt.Println(str)

	return str
}