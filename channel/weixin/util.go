package weixin

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func calculateSign(v WeixinRequest, md5Key string) (sign string) {
	dict := toMapWithValueNotNil(v)

	// eliminate any xml tag with value "-"
	delete(dict, "-")

	var keys []string
	for k, _ := range dict {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buffer bytes.Buffer
	for _, v := range keys {
		buffer.WriteString(v + "=" + dict[v] + "&")
	}
	buffer.WriteString("key=" + md5Key)

	seq := buffer.String()
	fmt.Println("seq:", seq)
	signSlice := md5.Sum([]byte(seq))

	return strings.ToUpper(hex.EncodeToString(signSlice[:]))
	// fmt.Println("sign:", microPay.Sign)
}

func toInt(s string) int {
	fmt.Println("s", s)
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return i
}

// map format: [device_info:123sdsf432dsf]
// only accept pointer value for now
func toMapWithValueNotNil(v interface{}) map[string]string {
	//req *micropayRequest
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	dict := make(map[string]string)
	if typ.Kind() == reflect.Ptr {
		for i := 0; i < typ.Elem().NumField(); i++ {
			tg := typ.Elem().Field(i).Tag.Get("xml")
			k := strings.Split(tg, ",")[0]
			f := val.Elem().Field(i)

			switch f.Kind() {
			case reflect.String:
				if f.Len() > 0 {
					v := f.String()
					dict[k] = v
				}
			case reflect.Int:
				v := int(f.Int())
				dict[k] = strconv.Itoa(v)

			case reflect.Struct:
				// do nothing
				// in case `XMLName xml.Name `xml:"xml"``

			default:
				panic("unsuported cast" + f.Kind().String())
			}
		}
		return dict
	} else {
		log.Fatal(errors.New("unsupported type: not pointer"))
	}
	return dict
}
