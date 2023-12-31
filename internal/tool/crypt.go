package tool

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"wizh/pkg/viper"
)

var (
	config = viper.InitConf("crypt")
	salt   = config.Viper.GetString("Salt")
	coder  = base64.NewEncoding(config.Viper.GetString("Base64Table"))
)

func Base64Encode(data []byte) []byte {
	return []byte(coder.EncodeToString(data))
}

func Base64Decode(data []byte) ([]byte, error) {
	return coder.DecodeString(string(data))
}

func Sha256Encrypt(data string) string {
	h := hmac.New(sha256.New, []byte(data+salt))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
