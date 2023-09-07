package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// Encode 编码
func Encode(text []byte, secretKey []byte) []byte {
	block, err := aes.NewCipher([]byte(fmt.Sprintf("%x", md5.Sum(secretKey))))
	if err != nil {
		return nil
	}
	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], text)
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(ciphertext)))
	base64.StdEncoding.Encode(dst, ciphertext)
	return dst
}

// Decode 解码
func Decode(ciphertext []byte, secretKey []byte) []byte {
	maxlen := base64.StdEncoding.DecodedLen(len(ciphertext))
	dst := make([]byte, maxlen)
	n, err := base64.StdEncoding.Decode(dst, ciphertext)
	if err != nil {
		return nil
	}
	ciphertext = dst[:n]
	block, err := aes.NewCipher([]byte(fmt.Sprintf("%x", md5.Sum(secretKey))))
	if err != nil {
		return nil
	}
	if len(ciphertext) < aes.BlockSize {
		return nil
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext
}
