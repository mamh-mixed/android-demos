package tools

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"github.com/omigo/g"
	"io/ioutil"
)

func EncodeBase64(data []byte) string {
	buff := &bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, buff)
	encoder.Write(data)
	encoder.Close()
	return buff.String()
}

func ToXML(data string) []byte {
	marshaled, err := xml.Marshal(data)
	if err != nil {
		g.Error("unable to marshal xml (%s)", err)
	}
	return marshaled
}

func EncodeHex(data []byte) string {
	return hex.EncodeToString(data)
}

func DecodeHex(data string) []byte {

	origin, err := hex.DecodeString(data)
	if err != nil {
		g.Error("fail to decode hex (%s)", err)
	}
	return origin
}

func DecodeBase64(data string) []byte {
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer([]byte(data)))
	origin, err := ioutil.ReadAll(decoder)
	if err != nil {
		g.Error("fail to decode base64 (%s)", err)
	}
	return origin
}
