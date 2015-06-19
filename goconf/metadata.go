package goconf

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

// 定义并验证所有配置项
// usage: goconf.Data.Mongo.URL

// SectionMongo Mongo 配置
type SectionMongo struct {
	URL string `meta:"url" fmt:"/.+/"`
	DB  string `meta:"db" fmt:"[a-zA-Z0-9]+" `
}

// SectionNSQ NSQ 配置
type SectionNSQ struct {
	// NSQLookupdAddr NSQ 服务发现地址
	NSQLookupdAddr string `meta:"nsqLookupd" fmt:"/[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}:[0-9]{4,5}/"`
	// 串码队列 Topic
	// ...
}

// Config 所有配置
type Config struct {
	Mongo SectionMongo `meta:"mongo"`
	NSQ   SectionNSQ   `meta:"nsq"`
}

// Data 配置数据
var Data = &Config{}

// tagMap 用来缓存从配置文件中读取的数据
var tagMap = make(map[string]map[string]string)

func metadata() {
	// fetch all tags of config struct
	fetchTags(reflect.TypeOf(Data).Elem())

	// validate pattern and value
	validate()

	// fill the correct value into map
	setTagMapValue()

	//initialize Data instance with correct value
	Data.initial()

	// show the result
	fmt.Println("Data:", Data)
}

// fetchTags 获得 Config 结构体中所有 field 的 meta 和 fmt
func fetchTags(t reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		tagMap[t.Field(i).Tag.Get("meta")] = fetchTagsFromSubFields(t.Field(i).Type)
	}
}

// fetchTagsFromSubFields
// subfields of Config contain `fmt` and `meta` tags
func fetchTagsFromSubFields(innerType reflect.Type) map[string]string {
	innerMap := make(map[string]string)
	for j := 0; j < innerType.NumField(); j++ {
		innerMap[innerType.Field(j).Tag.Get("meta")] = innerType.Field(j).Tag.Get("fmt")
	}
	return innerMap
}

// validate 校验 Config 每个 field 的格式
func validate() {
	for section, innerMap := range tagMap {
		for key, fmt := range innerMap {
			v := GetValue(section, key)
			if err := match(fmt, v); err != nil {
				panic(err)
			}
		}
	}
}

// setValue 将 ini 文件中的值存到 tagMap 里面
func setTagMapValue() {
	for section, innerMap := range tagMap {
		for key := range innerMap {
			innerMap[key] = GetValue(section, key)
		}
	}
}

// match 匹配正则表达式
func match(pattern, value string) error {
	match, err := regexp.MatchString(pattern, value)
	if err != nil {
		// some internal error occur
		return err
	}
	if !match {
		return errors.New("Error: " + value + " doesn't match " + pattern)
	}
	return nil
}

func initializeStruct(t reflect.Type, v reflect.Value, args ...interface{}) reflect.Value {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		ft := f.Type

		rawValue := ""
		// fetch tag
		tag := f.Tag.Get("meta")
		if len(args) == 0 {
			// means metaTag is section
			// do nothing ...
		} else {
			// fetch the section meta
			section := args[0].(string)
			key := tag
			rawValue = getValueFromTagMap(section, key)
		}

		ptr := reflect.New(ft)
		switch ft.Kind() {
		// only supprot int, string for now
		case reflect.Int:
			i, err := strconv.Atoi(rawValue)
			if err != nil {
				panic(err)
			}
			ptr.Elem().SetInt(int64(i))
			v.Elem().Field(i).Set(ptr.Elem())
		case reflect.String:
			ptr.Elem().SetString(rawValue)
			v.Elem().Field(i).Set(ptr.Elem())
		case reflect.Struct:
			ptr = initializeStruct(ft, ptr, tag)
			v.Elem().Field(i).Set(ptr.Elem())
		default:
		}
	}
	return v
}

// initial 初始化Config结构体
func (Data *Config) initial() {
	typeOfData := reflect.TypeOf(Data).Elem() // struct type
	ptrValueOfData := reflect.New(typeOfData)
	*Data = initializeStruct(typeOfData, ptrValueOfData).Elem().Interface().(Config)
}

// getValueFromTagMap 获取 map 中某一 meta 对应的 value
func getValueFromTagMap(section, key string) string {
	rawValue, ok := tagMap[section][key]
	if !ok {
		panic(errors.New("Errors: no value found in the map"))
	}
	return rawValue
}
