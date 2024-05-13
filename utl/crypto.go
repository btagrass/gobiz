package utl

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/spf13/cast"
)

var (
	Key = []byte("0123456789ABCDEFFEDCBA9876543210")
	iv  = []byte("0123456776543210")
)

func Decrypt(data string) (string, error) {
	inData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(Key)
	if err != nil {
		return "", err
	}
	cbc := cipher.NewCBCDecrypter(block, iv)
	outData := make([]byte, len(inData))
	cbc.CryptBlocks(outData, inData)
	outData = unPaddingPkcs5(outData)
	return string(outData), nil
}

func Digest(data string) string {
	md := md5.New()
	md.Write([]byte(data))
	bytes := md.Sum(nil)
	return base64.StdEncoding.EncodeToString(bytes)
}

func Encrypt(data string) (string, error) {
	block, err := aes.NewCipher(Key)
	if err != nil {
		return "", err
	}
	inData := paddingPkcs5([]byte(data), block.BlockSize())
	cbc := cipher.NewCBCEncrypter(block, iv)
	outData := make([]byte, len(inData))
	cbc.CryptBlocks(outData, inData)
	return base64.StdEncoding.EncodeToString(outData), nil
}

func HmacSha256(data string, secret string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(data))
	bytes := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(bytes)
}

func Md5(data string) string {
	bytes := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", bytes)
}

func paddingPkcs5(data []byte, blockSize int) []byte {
	paddingSize := blockSize - len(data)%blockSize
	paddings := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)
	return append(data, paddings...)
}

func unPaddingPkcs5(data []byte) []byte {
	length := len(data)
	unPadding := cast.ToInt(data[length-1])
	return data[:(length - unPadding)]
}
