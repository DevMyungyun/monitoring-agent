package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
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
