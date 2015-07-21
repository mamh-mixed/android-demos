package alipay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"github.com/omigo/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	AppId      = "2015032400038629"
	CusPrivKey = "MIICeQIBADANBgkqhkiG9w0BAQEFAASCAmMwggJfAgEAAoGBAL/74Ti5+Fv+0jRKXutpYnb+z9WMz0gADXK6ALNrVaLihqHucTdSNz6U7UCLWWwmMMmYTEvh0l1nkC+iSkyyfPh9PZCy6Y5gURIIhBufySztSJ+NsQZQHV0wMaVezjXtKKN9owwFIAySYgeVIzNyigOcMdPGlDi7cRbNasCqB0j9AgMBAAECgYEAn0b1n/6KUqii9McO4PhZvKfC/kUIY4/HmHtAwZo3Ph/52rbcLy2Cr+UWwQnbcqJsr5QvGHWN9fhJ43sdcWxTyNNwJecSjaWjta/JcWTByErHoNqJbj52EpiqQ+9PjzH5zBpdqKNx8YuaHL7/Y+7Ql8LuoDfXKo2TTcssZhj84PECQQD0En7PDX0E6a7B2BXiTJ+To5Q+OME9A+2XDESGswfHKW2/WbMpKM8LK/sU7JSwCRXa2jlTNkYCWiYu7QrYWhN7AkEAyV25lbMaJriP6cWS1A14Bz/uKk0U4YKO2TZrZvVjlwKqW5xoA8kwPlQCxxUG1zLozelbjg56wHt5gZJdQjqP5wJBAJbK3oG5yaXBYnDsugiIYobqp2oR0oGJ7b5GnAfEkGeh1uZD2wbw6YnzcDqrN+nSkygVbxlUDMjjPXf8h5jHfgUCQQDHJ2QyE3YMz8K90UMbaMrKWMdDnQLG2mpPmAv3Q0EhDGjSvEj/XY7SRiKNJVWjpt0rMd30DIwJLNWKeei0ZNkHAkEAxpQUui0qhg9g1HGVGpGgl1Jx5ULnj7AFYHwsFsm0nYXJIiJwdNKbC/B/GQ5wHnE22geYdxLS6L7zIQSAWMwTcQ=="
	ServerURL  = "https://openapi.alipay.com/gateway.do?charset=utf-8"
)

type PreOrder struct {
	OutTradeNo  string `json:"out_trade_no"`
	TotalAmount string `json:"total_amount"`
	Subject     string `json:"subject"`
}

func GenPostData(p *PreOrder) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		log.Error(err)
	}

	params := make(map[string]string)
	params["charset"] = "utf-8"
	params["sign_type"] = "RSA"
	params["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	params["method"] = "alipay.trade.precreate"
	params["biz_content"] = string(bytes)
	params["version"] = "1.0"
	params["app_id"] = AppId

	data := preContent(params)
	sign, _ := RsaSign(data, CusPrivKey)
	params["sign"] = sign

	val := url.Values{}
	for k, v := range params {
		val.Set(k, v)
	}

	return val.Encode()
}

func ProcessPreOrder() {
	p := &PreOrder{
		OutTradeNo:  "2015072017250000",
		TotalAmount: "0.01",
		Subject:     "讯联数据测试",
	}

	post(GenPostData(p))
}

func post(data string) {

	log.Debug(data)
	result, err := http.Post(ServerURL, "application/x-www-form-urlencoded;charset=utf-8", strings.NewReader(data))
	if err != nil {
		log.Error(err)
		return
	}
	defer result.Body.Close()

	resp, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debug(string(resp))

}

// Sign 签名
func RsaSign(content, cusPrivKey string) (string, error) {

	//to rsa.privateKey
	privKey := genPrivKeyFromPKSC8(cusPrivKey)
	hashed := sha1.Sum([]byte(content))
	signed, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA1, hashed[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signed), nil
}

func genPrivKeyFromPKSC8(pkcs8Key string) (privkey *rsa.PrivateKey) {

	// 解base64 pck8
	encodedKey, err := base64.StdEncoding.DecodeString(pkcs8Key)
	if err != nil {
		log.Error(err)
	}
	// 使用pkcs8格式
	pkcs8, err := x509.ParsePKCS8PrivateKey(encodedKey)

	var ok bool
	if privkey, ok = pkcs8.(*rsa.PrivateKey); !ok {
		log.Error(ok)
	}

	return
}
