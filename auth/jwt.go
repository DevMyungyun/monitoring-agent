package jwt

import (
	"net/http"
	"strings"
	"io/ioutil"
	"crypto/aes"
	"bytes"
	"encoding/json"
	"crypto/sha512"
	"encoding/hex"
	
	encryption "monitoring-agent/encryption"
	log "github.com/sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
)

func TokenValidCheck(headers *http.Request) interface{}{
	log.Info("=====token check=============")
	token := ExtractToken(headers)
	// log.Info("jwt token : ",token)
	// println(Cron.Entry(1))s

    data, err := ioutil.ReadFile("C:\\temp\\agent-config")
    if err != nil {
        panic(err)
    }

	key := "16byteSecret!!!!" // must be 16 byte
	block, err := aes.NewCipher([]byte(key))
	if err != nil {	
		log.Info(err)
	}
	// AES 알고리즘 암호문을 평문으로 복호화
	plaintext := encryption.Decrypt(block, data) 
	type agentData struct{
		ID string
		JWT string
	}
	var ad agentData

	editedPlaintext := bytes.Trim(plaintext, "\x00")

	error := json.Unmarshal(editedPlaintext, &ad)
	if error != nil {
		defer func() { 
			s := recover()
			log.Info(s) 
			log.Info(err)
		}()
        panic(error)
	}

	aStringToHash := []byte(ad.ID)
	sha512Bytes := sha512.Sum512(aStringToHash)
	secret := hex.EncodeToString(sha512Bytes[:])
	// log.Info("SHA512 String is ", hex.EncodeToString(sha512Bytes[:]))

	claim, check := VerifyToken(token, secret)
	if check != true {
		log.Info(check)
	}
	log.Info("result : ",claim)
	return claim
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	if bearToken != "" {
		strArr := strings.Split(bearToken, " ")
		if len(strArr) == 2 {
			return strArr[1]
		 }
		 return strArr[1]
	} else {
		log.Info("bearToken", bearToken)
		return ""
	}
}

func VerifyToken(myToken string, myKey string) (jwt.MapClaims, bool) {
	hmacSecretString := myKey
	hmacSecret := []byte(hmacSecretString)
	token, ok := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
			// check token signing method etc
			return hmacSecret, nil
	})
	if ok != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
	// token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
    //     return []byte(myKey), nil
    // })

    // if err == nil && token.Valid {
	// 	fmt.Println("Your token is valid.  I like your style.")
    // } else {
	// 	fmt.Println("This token is terrible!  I cannot accept this.")
	// }
	// return token, err
}

