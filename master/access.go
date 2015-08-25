package master

import (
	"fmt"
	"log"
	"net/http"

	"qiniupkg.com/api.v7/kodo"
)

const (
	BUCKET     = "test"
	DOMAIN     = "7xl02q.com1.z0.glb.clouddn.com"
	ACCESS_KEY = "-OOrgfZJbxz29kiW6HQsJ_OQJcjX6gaPRDf6xOcc"
	SECRET_KEY = "rgBxbGeGJluv8ApEjY1RL2vq9IIfXcQAQqH4ttGo"
)

var cli *kodo.Client

func init() {
	kodo.SetMac(ACCESS_KEY, SECRET_KEY)
	zone := 0                 // 您空间(Bucket)所在的区域
	cli = kodo.New(zone, nil) // 用默认配置创建 Client
}

func handleUptoken(w http.ResponseWriter, req *http.Request) {
	policy := &kodo.PutPolicy{
		Scope:   BUCKET,
		EndUser: "userId",
	}
	token := cli.MakeUptoken(policy)
	log.Println("token:", token)
	ret := fmt.Sprintf(`{"uptoken":"%s"}`, token)
	w.Write([]byte(ret))
}

func handleDownURL(w http.ResponseWriter, req *http.Request) {
	baseUrl := kodo.MakeBaseUrl(DOMAIN, req.URL.Query().Get("key")) // 得到下载 url
	img := cli.MakePrivateUrl(baseUrl, nil)                         // 用默认的下载策略去生成私有下载的 url

	// 如果是资质文件，需要保存路径或 key 值
	// 如果是 Excel/csv，需要下载并处理数据

	w.Write([]byte(img))
}
