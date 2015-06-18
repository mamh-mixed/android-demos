package weixin

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

func getRandomStr() string {
	return "sdfsfdsdfsfds"
}

// map format: [device_info:123sdsf432dsf]
// only accept pointer value for now
func toMapWithKeySortedAndValueNotNil(v interface{}) map[string]string {
	//req *micropayRequest
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	dict := make(map[string]string)
	if typ.Kind() == reflect.Ptr {
		// just one level deep

		for i := 0; i < typ.Elem().NumField(); i++ {
			k := typ.Elem().Field(i).Tag.Get("xml")
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
			default:
				panic("unsuported cast")
			}
		}
		return dict
	} else {
		log.Fatal(errors.New("unsupported type: not pointer"))
	}
	return dict
}
