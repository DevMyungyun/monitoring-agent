package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
)

func Encrypt(b cipher.Block, plaintext []byte) []byte {
	if mod := len(plaintext) % aes.BlockSize; mod != 0 { 
		padding := make([]byte, aes.BlockSize-mod)
		plaintext = append(plaintext, padding...)  
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext)) 
	iv := ciphertext[:aes.BlockSize] 
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println(err)
		return nil
	}
	mode := cipher.NewCBCEncrypter(b, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext) 
	return ciphertext
}

func Decrypt(b cipher.Block, ciphertext []byte) []byte {
	if len(ciphertext)%aes.BlockSize != 0 { 
		fmt.Println("Length of encrypted data must be multiple of size of block.")
		return nil
	}
	iv := ciphertext[:aes.BlockSize]        
	ciphertext = ciphertext[aes.BlockSize:] 
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(b, iv)                                                   
	mode.CryptBlocks(plaintext, ciphertext)   
	return plaintext
}

func GetDecryptData(secretKey string, data []byte) interface{} {
	key := secretKey // must be 16 byte
	block, err := aes.NewCipher([]byte(key))
	if err != nil {	
		log.Info(err)
	}
	// Decrypt with AES algorithm
	plaintext := Decrypt(block, data) 
	type agentData struct{
		ID string
		NAME string
		JWT string
		MAIN_SERVER_ADDRESS string
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

	config := make(map[string]string)
	config["ID"] = ad.ID
	config["NAME"] = ad.NAME
	config["JWT"] = ad.JWT
	config["MAIN_SERVER_ADDRESS"] = ad.MAIN_SERVER_ADDRESS
	log.Info("result : ",config)

	return config
}