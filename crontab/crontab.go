package crontab

import (
	"io/ioutil"
	"reflect"
	"encoding/json"

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

		res := httpReqRes.HttpReq("POST", ad["MAIN_SERVER_ADDRESS"]+urlPath, header, query, resource)
		log.Info("response : ", res)
		// log.Info("response type : ", reflect.TypeOf(res))
		
		 var resData map[string]string
		 json.Unmarshal([]byte(res),&resData)
		 log.Info("res json", resData)
		 log.Info("res json", resData["code"])
		 log.Info("res json type", reflect.TypeOf(resData))
		 // When JWT Token is expired
		 if resData["code"] == "401" {
			 log.Info("#### body : ", query)
			refreshTokenRes := httpReqRes.HttpReq("POST", ad["MAIN_SERVER_ADDRESS"]+"/auth/jwt/refresh/token", nil, query, query)
			log.Info("refreshToekn : ", refreshTokenRes)
		 }
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