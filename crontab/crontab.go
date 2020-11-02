package crontab

import (
	"io/ioutil"

	encryption "monitoring-agent/encryption"
	httpReqRes "monitoring-agent/http"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var Cron = cron.New()


func Intialize (timeSetting string, checkOS string, resource interface{}) {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	
	// Cron.WithSeconds()
	// Cron.AddFunc("*/1 * * * *", func() {
	Cron.AddFunc(timeSetting, func() {
		log.Info("[Job 1] ",timeSetting,"\n")
		// log.Info(">>>>", reflect.TypeOf(resource))
		// log.Info(">>>>", resource)
		
		filePath := ""
		urlPath := ""
		switch checkOS {
		case "windows":
			filePath = "C:\\temp\\agent-config"
			urlPath = "/agent/v1/windows/resource/receive"
		case "darwin":
			filePath = "/tmp/agent-config"
			urlPath = "/agent/v1/macos/resource/receive"
		case "linux":
			filePath = "/tmp/agent-config"
			urlPath = "/agent/v1/linux/resource/receive"
		default:
			log.Info("This OS is not supported : ", checkOS)
		}

		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		key := "16byteSecret!!!!" // must be 16 byte
		config := encryption.GetDecryptData(key, data)
		ad := config.(map[string]string)
		// Declare Header
		header := make(map[string]string)
		header["Content-Type"] = "application/json"
		header["Authorization"] = "bearer "+ ad["JWT"]
		// Declare Query
		query := make(map[string]string)
		query["name"] = ad["NAME"]

		url := ad["MAIN_SERVER_ADDRESS"]+urlPath
		res := httpReqRes.HttpReq("POST", url, header, query, resource)
		log.Info("response : ", res)
	})
}

func Start() {
	Cron.Start()
}

func Stop() {
	Cron.Stop()
}

func Entries() interface{} {
	return Cron.Entries()
}