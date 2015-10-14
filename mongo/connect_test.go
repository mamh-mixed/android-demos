package mongo

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/security"
	"github.com/omigo/log"
)

func TestEncryptMongoURL(t *testing.T) {
	url := goconf.Config.Mongo.URL
	t.Logf("origin: %s", url)

	ciphertext, err := security.RSAEncrypt([]byte(url), publicKey)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	b64Cipher := base64.StdEncoding.EncodeToString(ciphertext)
	t.Logf("encrypt: %s", b64Cipher)
}

func TestDecryptMongoURL(t *testing.T) {
	url := goconf.Config.Mongo.EncryptURL
	t.Logf("encrypt: %s", url)

	origData, err := security.RSADecryptBase64(url, privateKey)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("origin: %s", origData)
}

type TestModel struct {
	Timestamp time.Time
	Value     string
}

func TestConnect(t *testing.T) {
	for {
		c := database.C("testConnect")

		m := &TestModel{time.Now(), "test..."}
		m2 := &TestModel{time.Now().Add(time.Millisecond), "test2..."}

		err := c.Insert(m, m2)
		if err != nil {
			log.Errorf("insert error %s", err)
		}

		cond := bson.M{"timestamp": m.Timestamp}
		change := bson.M{"$set": bson.M{"updateTime": time.Now()}}
		err = c.Update(cond, change)
		if err != nil {
			log.Errorf("update error %s", err)
		}

		result := TestModel{}
		err = c.Find(cond).One(&result)
		if err != nil {
			log.Errorf("find error %s", err)
		}

		fmt.Println("timestamp:", result.Timestamp)

		err = c.Remove(cond)

		time.Sleep(time.Second)
	}
}
