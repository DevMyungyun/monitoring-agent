package jwt

import (
	"net/http"
	"io/ioutil"
	// "crypto/aes"
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

func HttpReq(method string, url string, header interface{}, query interface{}, body interface{}) string {
	bodyBytes, _ := json.Marshal(body)
	bodyBuffer := bytes.NewBuffer(bodyBytes)
	// log.Info("### ", bodyBytes)
	// log.Info("### ", bodyBuffer)
	
	// Generate Request Object
	req, err := http.NewRequest(method, url, bodyBuffer)
	if err !=nil {
		panic(err)
	}

	q := req.URL.Query()
	if query != nil {
		queries := query.(map[string]string)
		for key, val := range queries {
			q.Add(key, val)
		}
	}
	req.URL.RawQuery = q.Encode()

	if header != nil {
		headers := header.(map[string]string)
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	
	// Execute Request by Client Object
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Output result
	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes) 
	log.Info("response ",str)

	return str
}