package qiniu

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/goconf"
	"golang.org/x/net/context"
	"io"
	"log"
	"net/http"
	"qiniupkg.com/api.v7/kodo"
)

const (
	ACCESS_KEY = "-OOrgfZJbxz29kiW6HQsJ_OQJcjX6gaPRDf6xOcc"
	SECRET_KEY = "rgBxbGeGJluv8ApEjY1RL2vq9IIfXcQAQqH4ttGo"
)

var cli *kodo.Client
var BUCKET = goconf.Config.Qiniu.Bucket
var DOMAIN = goconf.Config.Qiniu.Domain

func init() {
	kodo.SetMac(ACCESS_KEY, SECRET_KEY)
	zone := 0                 // 您空间(Bucket)所在的区域
	cli = kodo.New(zone, nil) // 用默认配置创建 Client
}

func HandleUptoken(w http.ResponseWriter, req *http.Request) {
	policy := &kodo.PutPolicy{
		Scope:   BUCKET,
		EndUser: "userId",
	}
	token := cli.MakeUptoken(policy)
	log.Println("token:", token)
	ret := fmt.Sprintf(`{"uptoken":"%s"}`, token)
	w.Write([]byte(ret))
}

// List 列举资源
func List(prefix, marker string, limit int) ([]kodo.ListItem, string, error) {
	return cli.Bucket(BUCKET).List(context.Background(), prefix, marker, limit)
}

// Put 上传文件
func Put(key string, size int64, reader io.Reader) error {
	ctx := context.Background()
	return cli.Bucket(BUCKET).Put(ctx, nil, key, reader, size, &kodo.PutExtra{})
}

func MakePrivateUrl(key string) string {
	baseUrl := kodo.MakeBaseUrl(DOMAIN, key) // 得到下载 url
	return cli.MakePrivateUrl(baseUrl, nil)  // 用默认的下载策略去生成私有下载的 url
}

func HandleDownURL(w http.ResponseWriter, req *http.Request) {
	img := MakePrivateUrl(req.URL.Query().Get("key"))
	// 如果是资质文件，需要保存路径或 key 值
	// 如果是 Excel/csv，需要下载并处理数据
	w.Write([]byte(img))
}
